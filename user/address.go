package user

import (
	"errors"
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/helper"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddAddress(c *gin.Context) {
	var input models.InputAddress

	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if message, err := helper.ValidateAddress(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": message,
		})
		return
	}
	var user models.User
	if err := db.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cannot found user"})
		return
	}

	address := models.Address{
		UserID:       int(userID),
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		Country:      input.Country,
		City:         input.City,
		PostalCode:   input.PostalCode,
		Landmark:     input.Landmark,
	}

	if err := db.Db.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Address added successfully"})

}

func EditAddress(c *gin.Context) {
	var input models.InputAddress
	var address models.Address
	addressID := c.Param("id")
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID

	if err := db.Db.Where("user_id = ?", userID).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Address not found for this user id"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if message, err := helper.ValidateAddress(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": message,
		})
		return
	}
	updateaddress := models.Address{
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		Country:      input.Country,
		City:         input.City,
		PostalCode:   input.PostalCode,
		Landmark:     input.Landmark,
	}

	db.Db.Model(&models.Address{}).Where("address_id=?", addressID).Updates(&updateaddress)
	c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully"})
}

func DeleteAddress(c *gin.Context) {
	var address models.Address
	AddressID := c.Param("id")
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID

	if err := db.Db.Where("user_id = ?", userID).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Address not found for this user id"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	db.Db.Where("address_id=?", AddressID).Delete(&address)
	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}
