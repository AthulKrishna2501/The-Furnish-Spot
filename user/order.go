package user

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/AthulKrishna2501/The-Furniture-Spot/util"
	"github.com/gin-gonic/gin"

	"github.com/plutov/paypal/v4"

	"gorm.io/gorm"
)

func Orders(c *gin.Context) {
	var input models.OrderInput
	var address models.Address

	var coupon models.Coupon
	var cart []models.Cart

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	claims, _ := middleware.GetClaims(c)
	userID := claims.ID

	if err := db.Db.Where("user_id=? AND address_id=?", userID, input.AddressID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if err := db.Db.Table("carts").
		Select("carts.*, products.product_id, products.price").
		Joins("left join products on carts.product_id = products.product_id").
		Where("carts.user_id = ?", userID).
		Scan(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not fetch cart", "details": err.Error()})
		return
	}

	if len(cart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	var totalAmount float64
	var totalQuantity int
	var orderItems []models.OrderItem

	var totalDiscount float64

	for _, item := range cart {
		productID := item.ProductID
		product := models.Product{}

		if err := db.Db.First(&product, productID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found", "product_id": productID})
			return
		}

		if product.Quantity < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product", "product_id": product.ProductID})
			return
		}

		itemPrice := float64(item.Quantity) * product.Price

		var offer models.Offer
		var itemDiscount float64
		if err := db.Db.Where("product_id = ?", productID).First(&offer).Error; err == nil {
			itemDiscount = (float64(offer.OfferPercentage) / 100) * itemPrice
			itemPrice -= itemDiscount
		}

		totalDiscount += itemDiscount
		totalAmount += itemPrice
		totalQuantity += item.Quantity

		orderItem := models.OrderItem{
			ProductID: productID,
			Quantity:  item.Quantity,
			Price:     itemPrice,
		}
		orderItems = append(orderItems, orderItem)

		product.Quantity -= item.Quantity
		if err := db.Db.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}
	}

	var couponDiscount float64
	if input.CouponCode != "" {
		if err := db.Db.Where("coupon_code = ? AND is_active = ?", input.CouponCode, true).First(&coupon).Error; err == nil {
			if totalAmount >= float64(coupon.MinPurchaseAmount) {
				if coupon.DiscountType == "percentage" {
					couponDiscount = (coupon.DiscountAmount / 100) * totalAmount
				} else {
					couponDiscount = coupon.DiscountAmount
				}
				totalAmount -= couponDiscount
			} else {
				fmt.Println("Minimum purchase amount for coupon not met.")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Minimum purchase amount for coupon not met."})
				return
			}
		} else {

			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Coupon"})
			return
		}
	}
	totalDiscount += couponDiscount

	switch input.Method {
	case "Paypal":
		Total, err := util.ConvertINRtoUSD(totalAmount)
		if err != nil {
			log.Printf("Could not convert INR to USD: %v\n", err)
		}

		RoundedTotal := math.Round(Total*100) / 100
		client, err := NewPayPalClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize PayPal client"})
			return
		}

		approvalURL, payPalOrderID, err := CreatePayPalPayment(client, RoundedTotal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PayPal order"})
			return
		}

		tempOrder := models.TempOrder{
			OrderID:       payPalOrderID,
			UserID:        int(userID),
			CouponID:      coupon.CouponID,
			Quantity:      totalQuantity,
			Discount:      int(totalDiscount),
			Total:         totalAmount,
			Status:        "Pending",
			Method:        input.Method,
			PaymentStatus: "Pending",
			OrderDate:     time.Now(),
		}
		db.Db.Create(&tempOrder)

		c.JSON(http.StatusOK, gin.H{"approval_url": approvalURL})
		return

	case "COD":
		order, err := createOrder(userID, input, orderItems, totalAmount, totalQuantity, totalDiscount, coupon)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		orderResponse := responsemodels.OrderResponse{
			UserID:         int(userID),
			OrderID:        order.OrderID,
			Quantity:       totalQuantity,
			DiscountAmount: totalDiscount,
			Total:          totalAmount,
			Status:         order.Status,
			Method:         order.Method,
			PaymentStatus:  order.PaymentStatus,
			OrderDate:      order.OrderDate,
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": orderResponse})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

}

func ReturnOrder(c *gin.Context) {
	var order models.Order
	var input models.ReturnOrder

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Db.Where("order_id =?", input.OrderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if input.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please give a reason"})
		return
	}
	if order.Status != "Delivered" || order.Status == "Returned" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot return order"})
		return
	}
	order.Status = "Returned"

	db.Db.Save(&order)

	for _, item := range order.OrderItems {
		db.Db.Model(&models.Product{}).Where("product_id = ?", item.ProductID).Update("quantity", gorm.Expr("quantity + ?", item.Quantity))
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order returned successfully"})
}

func createOrder(userID uint, input models.OrderInput, orderItems []models.OrderItem, totalAmount float64,
	totalQuantity int, totalDiscount float64, coupon models.Coupon) (*models.Order, error) {

	order := models.Order{
		UserID:        int(userID),
		CouponID:      coupon.CouponID,
		Quantity:      totalQuantity,
		Discount:      int(totalDiscount),
		Total:         totalAmount,
		Status:        "Pending",
		Method:        input.Method,
		PaymentStatus: "Pending",
		OrderDate:     time.Now(),
	}

	if input.Method == "COD" {
		order.PaymentStatus = "Pending"
	} else if input.Method == "Paypal" {
		order.PaymentStatus = "Processing"
	}

	tx := db.Db.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("could not place order: %v", err)
	}

	for _, item := range orderItems {
		item.OrderID = order.OrderID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("could not save order items: %v", err)
		}
	}

	if err := tx.Where("user_id=?", userID).Delete(&models.Cart{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to clear cart: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &order, nil
}

func CapturePayPalOrder(c *gin.Context) {
	client, err := NewPayPalClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize PayPal client"})
		return
	}

	orderID := c.Query("token")
	payerID := c.Query("PayerID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID (token) missing from query parameters"})
		return
	}

	fmt.Printf("Received OrderID (token): %s, PayerID: %s\n", orderID, payerID)

	captureRequest := paypal.CaptureOrderRequest{}
	order, err := client.CaptureOrder(context.Background(), orderID, captureRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to capture PayPal order", "details": err.Error()})
		return
	}

	if order.Status != "COMPLETED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not completed"})
		return
	}

	var tempOrder models.TempOrder
	if err := db.Db.Where("order_id = ?", orderID).First(&tempOrder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Temporary order not found"})
		return
	}

	tempOrder.Status = "Processing"
	tempOrder.PaymentStatus = "Completed"
	if err := db.Db.Save(&tempOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}
	originalOrder := models.Order{
		PaymentID:     tempOrder.OrderID,
		UserID:        tempOrder.UserID,
		CouponID:      tempOrder.CouponID,
		Quantity:      tempOrder.Quantity,
		Discount:      tempOrder.Discount,
		Total:         tempOrder.Total,
		Status:        tempOrder.Status,
		Method:        tempOrder.Method,
		PaymentStatus: tempOrder.PaymentStatus,
		OrderDate:     tempOrder.OrderDate,
	}

	if err := db.Db.Create(&originalOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create original order"})
		return
	}
	if err := db.Db.Where("user_id=?", tempOrder.UserID).Delete(&models.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot clear cart"})
		return
	}

	orderResponse := gin.H{
		"order_id":       tempOrder.OrderID,
		"user_id":        tempOrder.UserID,
		"coupon_id":      tempOrder.CouponID,
		"quantity":       tempOrder.Quantity,
		"discount":       tempOrder.Discount,
		"total":          tempOrder.Total,
		"status":         tempOrder.Status,
		"method":         tempOrder.Method,
		"payment_status": tempOrder.PaymentStatus,
		"order_date":     tempOrder.OrderDate,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful and order created",
		"order":   orderResponse,
	})
}
