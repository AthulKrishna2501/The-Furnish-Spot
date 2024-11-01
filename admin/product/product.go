package product

import (
	"fmt"
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ViewProducts(c *gin.Context) {
	var products []models.Product
	result := db.Db.Order("product_id ASC").Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	for i := range products {
		if products[i].Quantity == 0 {
			products[i].Status = "Out of stock"
			db.Db.Model(&products[i]).Update("status", products[i].Status)
		}
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

	var category models.Category
	if err := db.Db.Where("category_id = ?", products.CategoryID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while checking category"})
		return
	}

	var existingProduct models.Product
	if err := db.Db.Where("product_name = ? AND category_id = ?", products.ProductName, products.CategoryID).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product already exists in this category"})
		return
	}
	if products.Price < 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Price cannot be a negative value"})
		return
	}
	p := models.Product{
		ProductName: products.ProductName,
		Price:       products.Price,
		CategoryID:  products.CategoryID,
		Description: products.Description,
		Quantity:    products.Quantity,
		Status:      products.Status,
		ImgURL:      products.ImgURL,
	}
	fmt.Println(p)
	if err := db.Db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
}

func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	var product models.Product

	if err := db.Db.Where("deleted_at IS NULL").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		ProductName string  `json:"product_name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImgURL      string  `json:"img_url"`
		Status      string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := models.Product{
		ProductName: input.ProductName,
		Description: input.Description,
		Price:       input.Price,
		ImgURL:      input.ImgURL,
		Status:      input.Status,
	}
	fmt.Println(updates)

	if err := db.Db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product_name": product.ProductName})
}

func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	var product models.Product

	if err := db.Db.Where("product_id = ?", productID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := db.Db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func UpdateProductStock(c *gin.Context) {
	var product models.Product
	productID := c.Param("id")

	if err := db.Db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Quantity = input.Quantity

	db.Db.Save(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Stock updated"})
}
