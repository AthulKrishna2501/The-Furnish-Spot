package user

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/AthulKrishna2501/The-Furniture-Spot/util"
	"github.com/gin-gonic/gin"
)

func Orders(c *gin.Context) {
	var input models.OrderInput
	var address models.Address
	var order models.Order
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

	var cartLock sync.Mutex
	cartLock.Lock()
	defer cartLock.Unlock()

	var totalAmount float64
	var totalQuantity int
	var orderItems []models.OrderItem

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
	var discount float64
	if input.CouponCode != "" {
		if err := db.Db.Where("coupon_code = ? AND is_active = ?", input.CouponCode, true).First(&coupon).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired coupon"})
			return
		}

		if totalAmount < float64(coupon.MinPurchaseAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order total does not meet coupon's minimum value requirement"})
			return
		}

		if coupon.DiscountType == "percentage" {
			discount = (coupon.DiscountAmount / 100) * totalAmount
		} else {
			discount = coupon.DiscountAmount
		}
		totalAmount -= discount
	}

	if input.Method == "Paypal" {
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

		approvalURL, err := CreatePayPalPayment(client, RoundedTotal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PayPal order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"approval_url": approvalURL})
		order.PaymentStatus = "Processing"

	} else if input.Method == "COD" {
		order.PaymentStatus = "Pending"

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	orders := models.Order{
		UserID:        int(userID),
		Total:         totalAmount,
		Quantity:      totalQuantity,
		Discount:      int(discount),
		CouponID:      coupon.CouponID,
		OrderDate:     time.Now(),
		Status:        "Pending",
		Method:        input.Method,
		PaymentStatus: order.PaymentStatus,
	}
	fmt.Println(orders)

	if err := db.Db.Create(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order", "details": err.Error()})
		return
	}

	for _, item := range orderItems {
		item.OrderID = orders.OrderID

		if err := db.Db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save order items", "details": err.Error()})
			return
		}
	}

	if err := db.Db.Where("user_id=?", userID).Delete(&models.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear the cart"})
		return
	}

	orderResponse := responsemodels.OrderResponse{
		UserID:         int(userID),
		OrderID:        orders.OrderID,
		Quantity:       totalQuantity,
		DiscountAmount: discount,
		Total:          totalAmount,
		Status:         orders.Status,
		Method:         orders.Method,
		PaymentStatus:  order.PaymentStatus,
		OrderDate:      orders.OrderDate,
	}
	fmt.Println(orderResponse)

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": orderResponse})

}
