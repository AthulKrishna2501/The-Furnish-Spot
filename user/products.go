package user

import (
	"fmt"
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ViewProducts(c *gin.Context) {
	var products []models.Product
	fmt.Println("HII")
	result := db.Db.Raw(` SELECT p.product_id, p.product_name, p.price, p.description,p.status,p.img_url,c.category_name,c.category_id FROM products p LEFT JOIN categories c 
   ON p.category_id = c.category_id WHERE  p.deleted_at IS NULL AND c.deleted_at IS NULL`).Scan(&products)
	fmt.Println(result)
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
