package category

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ViewCategory(c *gin.Context) {
	var category []models.Category

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

	if err := db.Db.Delete(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
