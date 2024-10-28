package user

import (
	"net/http"
	"sync"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	cartLock    sync.Mutex
	MaxQuantity = 5
)

func Cart(c *gin.Context) {
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID
	var cartItems []responsemodels.CartResponse

	if err := db.Db.Table("carts").
		Select("carts.user_id, carts.product_id, carts.quantity, carts.total, users.user_name, users.email").
		Joins("join users on users.id = carts.user_id").
		Where("carts.user_id = ?", userID).
		Scan(&cartItems).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	if len(cartItems) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Cart is empty"})
		return
	}
	c.JSON(http.StatusOK, cartItems)
}

func AddToCart(c *gin.Context) {
	var item models.CartInput
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID

	cartLock.Lock()
	defer cartLock.Unlock()

	var cartItem models.Cart
	var product models.Product

	if err := db.Db.First(&product, item.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if item.Quantity > product.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"message": "There is no sufficient quantity"})
		return
	}

	if item.Quantity > MaxQuantity {
		c.JSON(http.StatusOK, gin.H{"message": "Quantity limit exceeded"})
		return
	}

	if err := db.Db.Where("user_id = ? AND product_id = ?", userID, item.ProductID).First(&cartItem).Error; err == nil {
		cartItem.Quantity += item.Quantity
		cartItem.Total = cartItem.Quantity * int(product.Price)
		product.Quantity -= item.Quantity
		if err := db.Db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating cart item"})
			return
		}
		db.Db.Save(&product)

		c.JSON(http.StatusOK, gin.H{"message": "Item quantity updated"})
		return
	} else if err == gorm.ErrRecordNotFound {

		newCartItem := models.Cart{
			UserID:    int(userID),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Total:     item.Quantity * int(product.Price),
		}
		if err := db.Db.Create(&newCartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding item to cart"})
			return
		}

		product.Quantity -= item.Quantity
		db.Db.Save(&product)

		c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart"})
		product.Quantity -= item.Quantity
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
}

func RemoveItem(c *gin.Context) {
	ProductID := c.Param("id")
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID

	var cart models.Cart

	if err := db.Db.Where("user_id =? AND product_id =?", userID, ProductID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	var product models.Product
	if err := db.Db.First(&product, cart.ProductID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find product"})
		return
	}

	product.Quantity += cart.Quantity
	if err := db.Db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	if err := db.Db.Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed successfully"})

}
