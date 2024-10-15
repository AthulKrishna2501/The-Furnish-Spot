package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ListProducts(c *gin.Context) {
	var products models.Product

	if err := db.Db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrive products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
