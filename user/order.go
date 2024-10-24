package user

import (
	"net/http"
	"sync"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func Orders(c *gin.Context) {
	var input models.OrderInput
	var address models.Address
	var cart []models.Cart

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID
	if err := db.Db.Where("user_id=? AND address_id =?", userID, input.AddressID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if err := db.Db.Where("user_id=?", userID).Preload("Product").Find(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found for this userid"})
		return
	}

	var cartlock sync.Mutex

	cartlock.Lock()
	defer cartlock.Unlock()

	var totalamount float64
	var orderItems []models.OrderItem
	for _, item := range cart {
		product := item.Product
		if int(product.Quantity) < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product", "product_id": product.ProductID})
			return
		}
		itemPrice := float64(item.Quantity) * product.Price
		totalamount += itemPrice
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
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
	order := models.Order{
		UserID:     int(userID),
		Total:      int(totalamount),
		Status:     "Pending",
		Method:     "COD",
		OrderItems: orderItems,
	}

	if err := db.Db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order"})
		return
	}

	if err := db.Db.Where("user_id=?", userID).Delete(&models.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear the cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": order})

}
