package user

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/helper"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var input models.SignupInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received input: %+v\n", input)

	if message, err := helper.ValidateAll(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": message,
		})
		return
	}
	captchaID := input.CaptchaID
	captchaSolution := input.Captcha

	if !captcha.VerifyString(captchaID, captchaSolution) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Captcha"})
		return
	}

	var loginmethod string
	db.Db.Model(&models.UserLoginMethod{}).Where("user_login_method_email = ?", input.Email).Pluck("login_method", &loginmethod)
	if loginmethod == "Google Authentication" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "please log in through google authentication",
			"data":    gin.H{},
		})
		return
	}

	var count int64
	db.Db.Raw(`SELECT COUNT(*) FROM users where email = ?`, input.Email).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "already registered email",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	otp, err := helper.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating OTP"})
		return
	}

	newOtpRecord := models.OTP{
		Email:  input.Email,
		Code:   otp,
		Expiry: time.Now().Add(time.Minute * 5),
	}
	if err := db.Db.Create(&newOtpRecord).Error; err != nil {
		log.WithFields(log.Fields{
			"Email": input.Email,
			"error": err,
		}).Error("error creating otp")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go helper.SendEmail(input.Email, otp)

	user := models.TempUser{
		UserName:    input.UserName,
		Email:       input.Email,
		Password:    string(hashedPassword),
		PhoneNumber: input.PhoneNumber,
	}
	if err:=db.Db.Create(&user).Error;err!=nil{
		log.WithFields(log.Fields{
			"UserName":input.UserName,
			"error":err,
		}).Error("error creating user")
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP send successfully"})
}
