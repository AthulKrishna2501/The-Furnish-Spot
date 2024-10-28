package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
)

func ViewProducts(c *gin.Context) {
	var products []responsemodels.Products

	result := db.Db.Order("product_id ASC").Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	for i := range products {
		if products[i].Quantity == 0 {
			products[i].Status = "Out of stock"
			db.Db.Save(&products)
		} else {
			products[i].Status = "Available"
			db.Db.Save(&products)
		}
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No products listed"})
		return
	}

	c.JSON(http.StatusOK, products)
}
