package product

import (
	"fmt"
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ViewProducts(c *gin.Context) {

	var dbProducts []struct {
		ProductID     int     `gorm:"column:product_id"`
		ProductName   string  `gorm:"column:product_name"`
		Description   string  `gorm:"column:description"`
		Price         float64 `gorm:"column:price"`
		OfferDiscount float64 `gorm:"column:offer_discount"`
		CategoryID    uint    `gorm:"column:category_id"`
		ImgURL        string  `gorm:"column:img_url"`
		Status        string  `gorm:"column:status"`
		Quantity      int     `gorm:"column:quantity"`
		AverageRating float64 `gorm:"column:average_rating"`
		TotalReviews  int     `gorm:"column:total_reviews"`
	}

	result := db.Db.Raw(`
		 SELECT 
			  p.product_id,
			  p.product_name,
			  p.description,
			  p.price,
			  p.category_id,
			  p.img_url,
			  p.status,
			  p.quantity,
			  p.offer_discount,
			  COALESCE(AVG(r.rating), 0) as average_rating,
			  COUNT(DISTINCT r.review_rating_id) as total_reviews
		 FROM products p
		 LEFT JOIN review_ratings r ON r.product_id = p.product_id
		 WHERE p.deleted_at IS NULL
		 GROUP BY p.product_id
	`).Find(&dbProducts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(dbProducts) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No products listed"})
		return
	}

	responseProducts := make([]responsemodels.Products, len(dbProducts))

	for i, dbProduct := range dbProducts {
		status := "Available"
		if dbProduct.Quantity == 0 {
			status = "Out of stock"
		}

		var recentReviews []models.ReviewRating
		if result := db.Db.Table("review_ratings").
			Select("review_rating_id, user_id, rating, comment, created_at").
			Where("product_id = ?", dbProduct.ProductID).
			Order("created_at DESC").
			Limit(3).
			Find(&recentReviews); result.Error != nil {
			log.WithFields(log.Fields{
				"ProductID": dbProduct.ProductID,
				"error":     result.Error,
			}).Error("error retrieving reviewratings")
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		reviewResponses := make([]responsemodels.ReviewRating, len(recentReviews))
		for j, review := range recentReviews {
			reviewResponses[j] = responsemodels.ReviewRating{
				ReviewID:  review.ReviewRatingID,
				UserID:    review.UserID,
				Rating:    review.Rating,
				Comment:   review.Comment,
				CreatedAt: review.CreatedAt,
			}
		}

		responseProducts[i] = responsemodels.Products{
			ProductID:     dbProduct.ProductID,
			ProductName:   dbProduct.ProductName,
			Description:   dbProduct.Description,
			Price:         dbProduct.Price,
			OfferDiscount: dbProduct.OfferDiscount,
			CategoryID:    dbProduct.CategoryID,
			ImgURL:        dbProduct.ImgURL,
			Status:        status,
			Quantity:      dbProduct.Quantity,
			AverageRating: dbProduct.AverageRating,
			TotalReviews:  dbProduct.TotalReviews,
			RecentReviews: reviewResponses,
		}
	}

	c.JSON(http.StatusOK, responseProducts)
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

	if err := db.Db.Where("product_id=?", productID).Delete(&product).Error; err != nil {
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
