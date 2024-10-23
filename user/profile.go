package user

import (
	"errors"
	"fmt"
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
		return
	}

	userID := customClaims.ID
	var user responsemodels.User

	result := db.Db.Where("id=?", userID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
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

	exists := db.Db.Where("email=? AND id !=?", input.Email, userID).First(&user)
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
	var address []models.Address
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	userID := customClaims.ID
	fmt.Println(userID)

	result := db.Db.Where("user_id = ?", userID).Find(&address)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Address not found"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(address) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No address found"})
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

	result := db.Db.Where("user_id=?", userID).Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No orders found"})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No orders found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": orders})

}

func CancelOrders(c *gin.Context) {
	var orders models.Order
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID

	OrderID := c.Param("order_id")

	if err := db.Db.Where("order_id= ? AND user_id = ?", OrderID, userID).First(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Order not found or unauthorized"})
		return
	}

	if err := db.Db.First(&orders, OrderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Order not found"})
		return
	}

	if orders.Status == "Canceled" || orders.Status == "Delivered" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot cancel order"})
		return
	}
	orders.Status = "Canceled"

	if err := db.Db.Save(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save orderstatus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order canceled Successfully"})

}

func ForgotPassword(c *gin.Context) {
	var user models.User
	var input models.NewPassword

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

	if message, err := helper.ValidateAll(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}

	NewPassword := models.User{
		Password: input.NewPassword,
	}
	if err := db.Db.Model(&user).Where("id=?", userID).Updates(&NewPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})

}
