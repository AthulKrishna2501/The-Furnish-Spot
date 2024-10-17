package product

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ViewProducts(c *gin.Context) {
	var products []models.Product
	result := db.Db.Raw(`
        SELECT p.product_id, p.product_name, p.price, c.category_name AS category_name
        FROM products p
        LEFT JOIN categories c ON p.product_id = c.category_id`).Scan(&products)

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

func AddProducts(c *gin.Context) {
	var products models.Product

	if err := c.ShouldBind(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Db.Create(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Products added successfully"})
}

func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	var product models.Product

	if err := db.Db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Product updated successfully": product.ProductName})
}

func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	var product models.Product

	if err := db.Db.Delete(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
