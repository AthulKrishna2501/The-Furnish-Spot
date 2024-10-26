package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/helper"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	editUser := models.User{
		UserName:    input.UserName,
		Email:       input.Email,
		Password:    string(hashedPassword),
		PhoneNumber: input.PhoneNumber,
	}

	result := db.Db.Model(&models.User{}).Where("id = ?", userID).Updates(editUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})

}

func ViewAddress(c *gin.Context) {
	var address []responsemodels.Address
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	userID := customClaims.ID
	fmt.Println(userID)

	result := db.Db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&address)

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
	var orderItems []models.OrderItem

	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID
	OrderID := c.Param("id")

	// Convert OrderID to integer if needed
	orderID, err := strconv.Atoi(OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Check if the order exists and belongs to the user
	if err := db.Db.Where("order_id = ? AND user_id = ?", orderID, userID).First(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Order not found or unauthorized"})
		return
	}

	// Ensure the order can be canceled
	if orders.Status == "Canceled" || orders.Status == "Delivered" || orders.Status == "Failed" || orders.Status == "Shipped" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot cancel order"})
		return
	}

	// Retrieve all items associated with the order
	if err := db.Db.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Order items not found"})
		return
	}

	// Restore product quantities for each item in the order
	for _, item := range orderItems {
		var product models.Product

		if err := db.Db.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product not found"})
			return
		}

		// Update the product quantity by adding back the quantity from the canceled order
		product.Quantity += item.Quantity

		if err := db.Db.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update product quantity"})
			return
		}
	}

	// Update the order status to "Canceled"
	orders.Status = "Canceled"
	if err := db.Db.Save(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order canceled successfully"})
}

func ForgotPassword(c *gin.Context) {
	var input models.NewPassword
	var user models.User

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
	if err := db.Db.Model(&models.User{}).Where("id = ?", userID).Select("password").First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	NewPassword := string(hashedPassword)
	if err := db.Db.Model(&models.User{}).Where("id = ?", userID).Update("password", NewPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	fmt.Println(string(NewPassword))

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})

}
