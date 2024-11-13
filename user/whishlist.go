package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ViewWhishlist(c *gin.Context) {
	var whishlists []responsemodels.Wishlist

	claims, _ := middleware.GetClaims(c)

	userID := claims.ID

	if err := db.Db.Where("user_id=? AND deleted_at IS NULL", userID).Find(&whishlists).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": userID,
		}).Error("Failed to Retrive whishlist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot fetch whislist"})
		return
	}
	if len(whishlists) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Whishlist is empty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Wishlist retrieved successfully",
		"Wishlist": whishlists,
	})
}

func AddToWhishlist(c *gin.Context) {
	var whislist models.Wishlist
	var product models.Product
	var input struct {
		ProductID int `json:"product_id"`
	}

	claims, _ := middleware.GetClaims(c)
	userID := claims.ID

	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"UserID":    userID,
			"ProductID": input.ProductID,
		}).Error("Error in binding input")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Db.Where("product_id=?", input.ProductID).First(&product).Error; err != nil {
		log.WithFields(log.Fields{
			"ProductID": input.ProductID,
		}).Error("Can't find the product")
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.Quantity == 0 {
		log.WithFields(log.Fields{
			"ProductID": input.ProductID,
		}).Error("Insufficient quantity")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot add product, Insufficient quantity"})
		return
	}

	if err := db.Db.Where("product_id=? AND user_id=?", input.ProductID, userID).First(&whislist).Error; err == nil {
		log.WithFields(log.Fields{
			"ProductID": input.ProductID,
		}).Error("Product aldready in whislist")
		c.JSON(http.StatusFound, gin.H{"message": "Product aldready in whishlist"})
		return
	}

	whislistItem := models.Wishlist{
		ProductID:   input.ProductID,
		UserID:      int(userID),
		ProductName: product.ProductName,
		Price:       int(product.Price),
		Quantity:    product.Quantity,
	}

	if err := db.Db.Create(&whislistItem).Error; err != nil {
		log.WithFields(log.Fields{
			"ProductID": input.ProductID,
			"UserID":    userID,
		}).Error("Cannot create whishlist")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create whislist"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Item added to Whislist"})

}

func WishlistRemoveItem(c *gin.Context) {
	var whishlist models.Wishlist

	claims, _ := middleware.GetClaims(c)

	userID := claims.ID

	var input struct {
		ProductID int `json:"product_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"UserID":    userID,
			"ProductID": input.ProductID,
		}).Error("Invalid input")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Db.Where("product_id =? AND user_id =?", input.ProductID, userID).First(&whishlist).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":    userID,
			"ProductID": input.ProductID,
		}).Error("Item not found in wishlist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found in wishlist"})
		return
	}
	if err := db.Db.Where("user_id=? AND product_id =?", userID, input.ProductID).Delete(&whishlist).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":    userID,
			"ProductID": input.ProductID,
		}).Error("Cannot delete product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete Item"})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "Item removed successfully"})
}

func ClearWishlist(c *gin.Context) {
	var whishlist []models.Wishlist

	claims, _ := middleware.GetClaims(c)
	userID := claims.ID

	if err := db.Db.Where("user_id=?", userID).Find(&whishlist).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": userID,
		}).Error("Cannot find whislist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Wishlist not found"})
		return
	}

	if len(whishlist) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "There is nothing to clear"})
		return
	}

	if err := db.Db.Where("user_id=?", userID).Delete(&whishlist).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": userID,
		}).Error("Cannot clear wishlist")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot clear wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wishlist successfully cleared"})
}
