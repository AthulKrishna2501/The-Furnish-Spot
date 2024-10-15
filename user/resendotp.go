package user

import (
	"log"
	"net/http"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/helper"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ResendOTP(c *gin.Context) {
	Email := c.Param("email")

	otp, err := helper.GenerateOTP()
	if err != nil {
		log.Fatal("Error in generating otp")
	}
	go helper.SendEmail(Email, otp)
	newOtpRecord := models.OTP{
		Email:  Email,
		Code:   otp,
		Expiry: time.Now().Add(time.Minute * 5),
	}
	result := db.Db.Model(&models.OTP{}).Where("email = ?", Email).Updates(newOtpRecord)
	if result.Error != nil {
		log.Println("Error updating OTP record:", result.Error)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "OTP resend successfull"})
	}

}
