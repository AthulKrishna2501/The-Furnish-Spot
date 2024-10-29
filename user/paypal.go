package user

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v4"
	log "github.com/sirupsen/logrus"
)

func PaypalOrder(c *gin.Context) {
	ClientID := os.Getenv("CLIENT_ID")
	Secret := os.Getenv("SECRET")
	// isSandbox := true

	client, err := paypal.NewClient(ClientID, Secret, paypal.APIBaseSandBox)
	if err != nil {
		log.Fatalf("Failed to create paypal client :%v", err)

	}
	client.SetLog(os.Stdout)

	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}
	fmt.Println("PayPal client initialized successfully")

	claims, _ := c.Get("claims")

	customClaims, ok := claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var input models.OrderInput
	var address models.Address
	var cart []models.Cart

	userID := customClaims.ID

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Failed to bind JSON:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Db.Where("user_id=? AND address_id", userID, input.AddressID).First(&address).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID":    userID,
			"AddressID": input.AddressID,
		}).Error("Failed to retrive address")

		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if err := db.Db.Table("carts").
		Select("carts.cart_id,products.product_id,products.price,carts.quantity").
		Joins("INNER JOIN products ON carts.product_id = products.product_id").
		Where("carts.user_id=?", userID).Scan(&cart).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": userID,
		}).Error("Could not fetch cart")
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not fetch cart", "details": err.Error()})
		return
	}

	log.Println("Fetched cart contents:%+v\n", cart)

	if len(cart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return

	}

	var cartLock sync.Mutex
	cartLock.Lock()
	defer cartLock.Unlock()

}
