package user

import (
	"context"
	"fmt"

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
	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{
			"UserID":    userID,
			"AddressID": input.AddressID,
			"error":     err,
		}).Error("error querying address")
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if err := db.Db.Table("carts").
		Select("carts.*, products.product_id, products.price").
		Joins("left join products on carts.product_id = products.product_id").
		Where("carts.user_id = ?", userID).
		Scan(&cart).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": userID,
			"error":  err,
		}).Error("error querying carts")
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
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error saving product")
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
			Discount:      totalDiscount,
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
			log.WithFields(log.Fields{
				"UserID": userID,
				"error":  err,
			}).Error("error creating order")
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
	case "Wallet":
		var wallet models.Wallet

		if err := db.Db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			log.WithFields(log.Fields{"UserID": userID}).Error("Cannot find wallet")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot find wallet"})
			return
		}

		if totalAmount > wallet.Balance {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Insufficient wallet balance"})
			return
		}

		order := models.Order{
			UserID:        int(userID),
			CouponID:      coupon.CouponID,
			Quantity:      totalQuantity,
			Discount:      totalDiscount,
			Total:         totalAmount,
			Status:        "Pending",
			Method:        input.Method,
			PaymentStatus: "Pending",
			OrderDate:     time.Now(),
		}

		tx := db.Db.Begin()

		wallet.Balance -= totalAmount
		if err := tx.Save(&wallet).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
			return
		}

		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		for _, item := range orderItems {
			item.OrderID = order.OrderID
			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add order item"})
				return
			}
		}

		if err := tx.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
			return
		}

		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
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

	if err := db.Db.Save(&order).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error saving order")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, item := range order.OrderItems {
		db.Db.Model(&models.Product{}).Where("product_id = ?", item.ProductID).Update("quantity", gorm.Expr("quantity + ?", item.Quantity))
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order returned successfully"})
}

func createOrder(userID uint, input models.OrderInput, orderItems []models.OrderItem, totalAmount float64,
	totalQuantity int, totalDiscount float64, coupon models.Coupon) (*models.Order, error) {

	if totalAmount > 1000 {
		return nil, fmt.Errorf("COD not allowed over Rupees 1000")
	}
	order := models.Order{
		UserID:        int(userID),
		CouponID:      coupon.CouponID,
		Quantity:      totalQuantity,
		Discount:      totalDiscount,
		Total:         totalAmount,
		Status:        "Pending",
		Method:        input.Method,
		PaymentStatus: "Pending",
		OrderDate:     time.Now(),
	}
	fmt.Println(order)

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
		log.WithFields(log.Fields{
			"PaymentID": originalOrder.PaymentID,
			"error":     err,
		}).Error("error creating order")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create original order"})
		return
	}

	var cartItems []models.Cart
	if err := db.Db.Where("user_id = ?", tempOrder.UserID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	var orderItems []models.OrderItem
	for _, cartItem := range cartItems {
		orderItem := models.OrderItem{
			OrderID:   originalOrder.OrderID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     (float64(cartItem.Total) / float64(cartItem.Quantity)),
		}
		orderItems = append(orderItems, orderItem)
	}

	if err := db.Db.Create(&orderItems).Error; err != nil {
		log.WithFields(log.Fields{
			"OrderID": originalOrder.OrderID,
			"error":   err,
		}).Error("error creating orderItems")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order items"})
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
