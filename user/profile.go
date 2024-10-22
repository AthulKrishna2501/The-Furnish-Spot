package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/helper"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserProfile(c *gin.Context) {
	claims, _ := c.Get("claims")

	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invlaid claims"})
	}

	userID := customClaims.ID
	var user responsemodels.User

	result := db.Db.Where("id=?", userID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"User Retrieved Successfully": user})

}

func EditProfile(c *gin.Context) {
	var user models.User

	claims, _ := c.Get("claims")

	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}
	userID := customClaims.ID

	var input models.EditUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists := db.Db.Where("email=?", input.Email).First(&user)
	if exists.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusConflict, gin.H{"message": "Email aldready exists"})
		return
	}

	message, err := helper.ValidateAll(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}

	editUser := models.User{
		UserName:    input.UserName,
		Email:       input.Email,
		Password:    input.Password,
		PhoneNumber: input.PhoneNumber,
	}

	result := db.Db.Model(&models.User{}).Where("id = ?", userID).Updates(editUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"User updated successfully": editUser})

}

func ViewAddress(c *gin.Context) {
	var address models.Address
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	userID := customClaims.ID

	result := db.Db.Where("user_id=?", userID).First(address)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": address})
}

func ViewOrders(c *gin.Context) {
	var orders []models.Order
	claims, _ := c.Get("claims")

	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	userID := customClaims.ID

	result := db.Db.Where("user_id=?", userID).First(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No orders found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": orders})

}
