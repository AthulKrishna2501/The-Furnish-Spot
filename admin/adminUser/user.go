package adminuser

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	var users []models.User

	if err := db.Db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Users": users})
}

func BlockUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	if err := db.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Status = "Blocked"
	if err := db.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

func UnblockUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	if err := db.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Status = "Active"
	if err := db.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Unblocked successfully"})
}
