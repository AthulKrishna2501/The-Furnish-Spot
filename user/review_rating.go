package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func AddReviews(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)

	UserID := claims.ID

	var input models.ReviewInput
	var product models.Product
	var review models.ReviewRating
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := db.Db.Where("user_id=? AND product_id =?", UserID, input.ProductID).First(&review).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Aldredy Reviewed for this product"})
		return
	}

	if err := db.Db.Where("product_id=?", input.ProductID).First(&product).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":    UserID,
			"ProductID": input.ProductID,
		}).Error("Product not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if input.Rating > 5 || input.Rating <= 0 {
		log.WithFields(log.Fields{
			"UserID":    UserID,
			"ProductID": input.ProductID,
		}).Error("Invalid rating")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid rating"})
		return
	}

	if input.Comment == "" {
		log.WithFields(log.Fields{
			"UserID":    UserID,
			"ProductID": input.ProductID,
		}).Error("Please give a comment")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please give a comment"})
		return
	}

	NewReview := models.ReviewRating{
		UserID:    int(UserID),
		ProductID: input.ProductID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	if err := db.Db.Create(&NewReview).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":    UserID,
			"ProductID": input.ProductID,
		}).Error("Cannot create review")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create review"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Review added succesfully"})
}

func EditReview(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)
	UserID := claims.ID

	var product models.Product
	var review models.ReviewRating
	var input models.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Db.Where("user_id=? AND product_id =?", UserID, input.ProductID).First(&review).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":    UserID,
			"ProductID": input.ProductID,
		}).Error("Cannot find Review")
		c.JSON(http.StatusNotFound, gin.H{"message": "Review not found"})
		return
	}
	if err := db.Db.Where("product_id=?", input.ProductID).First(&product).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":    UserID,
			"ProductID": input.ProductID,
		}).Error("Product not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if input.Rating > 5 || input.Rating <= 0 {
		log.WithFields(log.Fields{
			"UserID":    UserID,
			"ProductID": input.ProductID,
		}).Error("Invalid rating")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid rating"})
		return
	}

	if input.Comment == "" {
		log.WithFields(log.Fields{
			"UserID":    UserID,
			"ProductID": input.ProductID,
		}).Error("Please give a comment")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please give a comment"})
		return
	}

	NewReview := models.ReviewRating{
		UserID:    int(UserID),
		ProductID: input.ProductID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	if err := db.Db.Where("user_id =? AND product_id =?", UserID, input.ProductID).Updates(&NewReview).Error; err != nil {
		log.WithFields(log.Fields{
			"ProductID": input.ProductID,
			"UserID":    UserID,
		}).Error("Cannot update review")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review udpated successfully"})

}

func DeleteReview(c *gin.Context) {
	claims, _ := middleware.GetClaims(c)

	UserID := claims.ID

	ReviewID := c.Param("id")

	var review models.ReviewRating

	if err := db.Db.Where("user_id=? AND review_rating_id =?", UserID, ReviewID).First(&review).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":   UserID,
			"ReviewID": ReviewID,
		}).Error("Review not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found "})
		return
	}
	if err := db.Db.Where("review_rating_id =?", ReviewID).Delete(&review).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":   UserID,
			"ReviewID": ReviewID,
		}).Error("Cannot delete review")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete review"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
