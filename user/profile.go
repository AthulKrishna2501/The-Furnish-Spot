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

	log "github.com/sirupsen/logrus"
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

	if user.Status == "" {
		user.Status = "Active"
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
		log.WithFields(log.Fields{
			"UserName": input.UserName,
			"error":    err,
		}).Error("error updating user")
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
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot find userID",
			"err":   result.Error,
		})
		return
	}

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

	result := db.Db.Where("user_id = ?", userID).Preload("OrderItems").Find(&orders)
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
	var wallet models.Wallet

	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID
	OrderID := c.Param("id")

	orderID, err := strconv.Atoi(OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := db.Db.Where("order_id = ? AND user_id = ?", orderID, userID).First(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Order not found or unauthorized"})
		return
	}

	if orders.Status == "Canceled" || orders.Status == "Delivered" || orders.Status == "Failed" || orders.Status == "Shipped" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot cancel order"})
		return
	}

	if err := db.Db.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Order items not found"})
		return
	}

	for _, item := range orderItems {
		var product models.Product

		if err := db.Db.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product not found"})
			return
		}

		product.Quantity += item.Quantity

		if err := db.Db.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update product quantity"})
			return
		}
	}

	if orders.Method == "Paypal" {

		if err := db.Db.Where("user_id=?", userID).First(&wallet).Error; err == nil {
			wallet.Balance += orders.Total
			if err := db.Db.Save(&wallet).Error; err != nil {
				log.WithFields(log.Fields{
					"UserID": userID,
				}).Error("Cannot update wallet")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating wallet"})
				return
			}
			walletTransaction := models.WalletTransaction{
				UserID:          userID,
				OrderID:         uint(orders.OrderID),
				Amount:          orders.Total,
				TransactionType: "Credit",
				Description:     "Refund for Order #" + strconv.Itoa(int(orders.OrderID)),
			}
			if err := db.Db.Create(&walletTransaction).Error; err != nil {
				log.WithFields(log.Fields{
					"UserID":  userID,
					"OrderID": orders.OrderID,
				}).Error("Cannot create transcation")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create Transaction"})
				return
			}

		}
	}
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
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	if input.NewPassword != input.ReEnter {
		log.WithFields(log.Fields{
			"UserID": userID,
		}).Error("Password mismatch")

		c.JSON(http.StatusBadRequest, gin.H{"message": "Password does not match"})
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

func ViewWallet(c *gin.Context) {
	var wallet models.Wallet
	claims, _ := middleware.GetClaims(c)
	userID := claims.ID

	if err := db.Db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			wallet = models.Wallet{
				UserID:  userID,
				Balance: 0,
			}
			if createErr := db.Db.Create(&wallet).Error; createErr != nil {
				log.WithFields(log.Fields{
					"UserID": userID,
				}).Error("Failed to create wallet")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
				return
			}
			log.WithFields(log.Fields{
				"UserID": userID,
			}).Info("New wallet created")
		} else {
			log.WithFields(log.Fields{
				"UserID":   userID,
				"WalletID": wallet.WalletID,
			}).Error("Error finding wallet")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding wallet"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"Balance": wallet.Balance})
}

func GetWalletTransactions(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	userID := claims.ID

	var transactions []models.WalletTransaction

	if err := db.Db.Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving transactions"})
		return
	}
	if len(transactions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No transactions found"})
		return
	}

	var ResTransactions []responsemodels.Transaction

	for _, transaction := range transactions {
		ResTransaction := responsemodels.Transaction{
			OrderID:         transaction.OrderID,
			Amount:          transaction.Amount,
			TransactionType: transaction.TransactionType,
			Description:     transaction.Description,
		}
		ResTransactions = append(ResTransactions, ResTransaction)
	}

	c.JSON(http.StatusOK, gin.H{
		"Transactions": ResTransactions,
	})
}
