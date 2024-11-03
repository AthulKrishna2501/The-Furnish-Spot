package offer

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ViewOffers(c *gin.Context) {
	var offers []responsemodels.Offer

	if err := db.Db.Find(&offers).Error; err != nil {
		log.WithFields(log.Fields{
			"ERROR": "Cannot retrieve offers",
		}).Error("Cannot retrieve offers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot retrieve offers"})
		return
	}

	if len(offers) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "There are no offers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Offers retrieved successfully": offers})
}

func AddOffer(c *gin.Context) {
	var input models.OfferInput
	var product models.Product

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Db.Where("product_id = ?", input.ProductID).First(&product).Error; err != nil {
		log.WithFields(log.Fields{
			"ProductID": input.ProductID,
		}).Error("Cannot find product")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot find product"})
		return
	}

	if input.OfferPercentage == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please enter an offer amount"})
		return
	}

	if input.OfferPercentage <= 0 || input.OfferPercentage > 100 {
		log.WithFields(log.Fields{
			"ProductID":   input.ProductID,
			"OfferAmount": input.OfferPercentage,
		}).Error("Invalid offer amount")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid offer amount"})
		return
	}

	NewOffer := models.Offer{
		ProductID:       input.ProductID,
		OfferPercentage: input.OfferPercentage,
	}

	db.Db.Create(&NewOffer)

	if err := db.Db.Model(&product).Where("product_id = ?", input.ProductID).
		Update("offer_discount", input.OfferPercentage).Error; err != nil {
		log.WithFields(log.Fields{
			"ProductID":   input.ProductID,
			"OfferAmount": input.OfferPercentage,
		}).Error("Cannot update offer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update offer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Offer updated successfully"})
}

func UpdateOffer(c *gin.Context) {
	var offer models.Offer
	var input models.OfferInput
	var product models.Product

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Db.Where("product_id = ?", input.ProductID).First(&product).Error; err != nil {
		log.WithFields(log.Fields{
			"ProductID": input.ProductID,
		}).Error("Cannot find product")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot find product"})
		return
	}

	if input.OfferPercentage < 0 || input.OfferPercentage > 100 {
		log.WithFields(log.Fields{
			"ProductID":   input.ProductID,
			"OfferAmount": input.OfferPercentage,
		}).Error("Invalid offer amount")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid offer amount"})
		return
	}

	if err := db.Db.Model(&offer).Where("product_id = ?", input.ProductID).
		Update("offer_percentage", input.OfferPercentage).Error; err != nil {
		log.WithFields(log.Fields{
			"ProductID":   input.ProductID,
			"OfferAmount": input.OfferPercentage,
		}).Error("Cannot update offer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update offer"})
		return
	}

	if err := db.Db.Model(&product).Where("product_id = ?", input.ProductID).
		Update("offer_discount", input.OfferPercentage).Error; err != nil {
		log.WithFields(log.Fields{
			"ProductID":   input.ProductID,
			"OfferAmount": input.OfferPercentage,
		}).Error("Cannot update offer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update offer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Offer updated successfully"})

}
