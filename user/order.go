package user

import (
	"net/http"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			}
		}
	}
	totalDiscount += couponDiscount

	orders := models.Order{
		UserID:        int(userID),
		CouponID:      coupon.CouponID,
		Quantity:      totalQuantity,
		Discount:      int(totalDiscount),
		Total:         totalAmount,
		Status:        "Pending",
		Method:        input.Method,
		PaymentStatus: order.PaymentStatus,
		OrderDate:     time.Now(),
	}

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
		DiscountAmount: totalDiscount,
		Total:          totalAmount,
		Status:         orders.Status,
		Method:         orders.Method,
		PaymentStatus:  order.PaymentStatus,
		OrderDate:      orders.OrderDate,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": orderResponse})

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
