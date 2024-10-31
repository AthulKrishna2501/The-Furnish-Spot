package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ViewCoupons(c *gin.Context) {
	var coupons []models.Coupon

	if err := db.Db.Find(&coupons).Error; err != nil {
		log.WithFields(log.Fields{
			"Error": "Cannot get coupos",
		}).Error("Error retriving coupons")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error retiving coupons"})
		return
	}

	if len(coupons) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No coupons listed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Coupons retrieved successfully": coupons})
}
