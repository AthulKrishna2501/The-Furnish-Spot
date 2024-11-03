package coupon

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ViewCoupons(c *gin.Context) {
	var coupons []models.Coupon

	if err := db.Db.Where("deleted_at IS NULL").Find(&coupons).Error; err != nil {
		log.WithFields(log.Fields{
			"Coupon": "Cannot retrieve coupons",
		}).Error("Cannot show coupons")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot find coupons"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": coupons})
}

func AddCoupon(c *gin.Context) {
	var coupon models.Coupon
	var input models.CouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Db.Where("coupon_code=?", input.CouponCode).First(&coupon).Error; err == nil {
		log.WithFields(log.Fields{
			"CouponID": coupon.CouponID,
		}).Error("Coupon aldready exists")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Coupon aldready exists"})
		return
	}

	if input.MinPurchaseAmount == 0 && input.MaxPurchaseAmount == 0 {
		log.WithFields(log.Fields{
			"error": "Cannot give both as zero",
		}).Error("Both are zero")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot give both as zero"})
		return
	}

	if input.DiscountType != "percent" && input.DiscountType != "fixed" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discount type"})
		return
	}

	if input.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please add a description"})
		return
	}

	NewCoupon := models.Coupon{
		CouponID:          coupon.CouponID,
		CouponCode:        input.CouponCode,
		DiscountAmount:    input.DiscountAmount,
		DiscountType:      input.DiscountType,
		Description:       input.Description,
		StartDate:         input.StartDate,
		EndDate:           input.EndDate,
		MaxPurchaseAmount: input.MaxPurchaseAmount,
		MinPurchaseAmount: input.MinPurchaseAmount,
		IsActive:          input.IsActive,
	}

	if err := db.Db.Create(&NewCoupon).Error; err != nil {
		log.WithFields(log.Fields{
			"CouponID": coupon.CouponID,
		}).Error("Cannot create coupon")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create coupon"})
		return

	}

	c.JSON(http.StatusCreated, gin.H{"message": "Coupon created successfully"})

}

func DeleteCoupon(c *gin.Context) {
	var coupon models.Coupon
	couponID := c.Param("id")

	if err := db.Db.Where("coupon_id=?", couponID).First(&coupon).Error; err != nil {
		log.WithFields(log.Fields{
			"CouponID": couponID,
		}).Error("Coupon not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot find coupon"})
		return
	}
	if err := db.Db.Where("coupon_id=?", couponID).Delete(&coupon).Error; err != nil {
		log.WithFields(log.Fields{
			"CouponID": couponID,
		}).Error("Cannot delete coupon")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete coupon"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Coupon deleted sucessfully"})

}
