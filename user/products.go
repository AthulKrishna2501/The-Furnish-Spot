package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
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
		db.Db.Table("review_ratings").
			Select("review_rating_id, user_id, rating, comment, created_at").
			Where("product_id = ?", dbProduct.ProductID).
			Order("created_at DESC").
			Limit(3).
			Find(&recentReviews)

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
