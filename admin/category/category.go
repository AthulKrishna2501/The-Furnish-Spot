package category

import (
	"net/http"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ViewCategory(c *gin.Context) {
	var category []models.Category
	result := db.Db.Raw(`
        SELECT category_id, category_name 
        FROM categories 
        WHERE deleted_at IS NULL`).Scan(&category)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if len(category) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No categories listed"})
		return
	}
	if err := db.Db.Find(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Categories": category})
}

func AddCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if category.CategoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var existingCategory models.Category
	if err := db.Db.Where("category_name = ?", category.CategoryName).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
		return
	}
	if err := db.Db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Category created successfully": category.CategoryName})
}

func EditCategory(c *gin.Context) {
	categoryID := c.Param("id")
	var category models.Category

	if err := db.Db.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input "})
		return
	}
	if err := db.Db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Category updated successfully": category.CategoryName})
}

func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var category models.Category

	if err := db.Db.Where("category_id = ?", categoryID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := db.Db.Model(&category).Update("deleted_at", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	if err := db.Db.Model(&models.Product{}).Where("category_id = ?", categoryID).Update("deleted_at", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
