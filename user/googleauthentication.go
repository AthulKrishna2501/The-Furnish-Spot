package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type GoogleResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

var googleOauthConfig *oauth2.Config
var oauthStateString = "randomstring"

func HandleGoogleLogin(c *gin.Context) {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://thefurnishstore.shop/auth/google/callback",
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	fmt.Println(googleOauthConfig.ClientID, googleOauthConfig.ClientSecret)
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGoogleCallback(c *gin.Context) {

	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid oauth state"})
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}

	userInfoResp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}
	defer userInfoResp.Body.Close()

	userInfo, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error": "Failed to read user info: " + err.Error()})
		return
	}

	var user GoogleResponse
	if err := json.Unmarshal(userInfo, &user); err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error": "Failed to unmarshal user info: " + err.Error()})
		return
	}

	var loginMethod string
	err = db.Db.Model(&models.UserLoginMethod{}).
		Select("login_method").
		Where("user_login_method_email = ?", user.Email).
		Scan(&loginMethod).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	if loginMethod == "Manual" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Please log in through Manual Log in",
			"data":    gin.H{},
		})
		return
	}

	if loginMethod == "" {
		userLoginMethod := models.UserLoginMethod{
			UserLoginMethodEmail: user.Email,
			LoginMethod:          "Google Authentication",
		}
		newUser := models.User{
			UserName: user.Name,
			Email:    user.Email,
		}

		if err := db.Db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&userLoginMethod).Error; err != nil {
				return err
			}
			return tx.Create(&newUser).Error
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}
	}

	var userID uint
	err = db.Db.Model(&models.User{}).
		Where("email = ?", user.Email).
		Pluck("id", &userID).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID: " + err.Error()})
		return
	}

	jwtToken, err := middleware.CreateToken("user", user.Email, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetCookie("jwt_token", jwtToken, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User login successful",
		"token":   jwtToken,
	})

}
