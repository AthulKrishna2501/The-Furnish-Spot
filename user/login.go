package user

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var input models.LoginInput
	var user models.User

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Db.Where("email=?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}
	token, err := middleware.CreateToken(user.UserName, user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Error Generating jwt"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successfull", "token": token})
}
