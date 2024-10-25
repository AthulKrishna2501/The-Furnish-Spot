package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
)

func ViewProducts(c *gin.Context) {
	var products []responsemodels.Product

	result := db.Db.Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No products listed"})
		return
	}

	c.JSON(http.StatusOK, products)
}
