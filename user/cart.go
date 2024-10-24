package user

import (
	"net/http"
	"sync"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
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
	var cart models.Cart
	if err := db.Db.Preload("User").Preload("Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, cart)

}

func AddToCart(c *gin.Context) {
	var item models.Cart
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

	var cart models.Cart

	if err := db.Db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			cart = models.Cart{UserID: int(userID)}
			if err := db.Db.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while creating cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	var product models.Product
	if err := db.Db.First(&product, item.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if item.Quantity > int(product.Quantity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requested quantity exceeds available stock"})
		return
	}

	var existingCartItem models.Cart
	if err := db.Db.Where("cart_id = ? AND product_id = ?", cart.CartID, item.ProductID).First(&existingCartItem).Error; err == nil {
		existingCartItem.Quantity += item.Quantity
		existingCartItem.Total = existingCartItem.Quantity * int(product.Price)
		if err := db.Db.Save(&existingCartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating cart item"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Item quantity updated"})
	} else if err == gorm.ErrRecordNotFound {
		if item.Quantity > MaxQuantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot exceed maximum quantity limit per user"})
			return
		}

		item.CartID = cart.CartID
		item.Total = item.Quantity * int(product.Price)
		if err := db.Db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding item to cart"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
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
