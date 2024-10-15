package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func SearchProduct(c *gin.Context) {
	var input models.SearchProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var products []models.Product

	if err := db.Db.Where("product_name ILIKE ?", "%"+input.Name+"%").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching product"})
		return
	}
	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No Products found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Produts": products})
}
