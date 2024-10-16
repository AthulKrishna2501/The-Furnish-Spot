package usermanagement

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ListUser(c *gin.Context) {
	var user models.User

	if err := db.Db.Find(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrive users"})
		return
	}
	c.JSON(http.StatusOK, user)
}
