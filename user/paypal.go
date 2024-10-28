package user

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v4"
)

func Paypal(c *gin.Context) {
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

}
