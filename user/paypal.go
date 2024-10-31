package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v4"
)

func NewPayPalClient() (*paypal.Client, error) {
	client, err := paypal.NewClient(os.Getenv("CLIENT_ID"), os.Getenv("SECRET"), paypal.APIBaseSandBox)
	if err != nil {
		return nil, err
	}
	client.SetLog(log.Writer())
	return client, nil
}

func CreatePayPalPayment(client *paypal.Client, amount float64) (string, error) {
	purchaseUnit := paypal.PurchaseUnitRequest{
		Amount: &paypal.PurchaseUnitAmount{
			Currency: "USD",
			Value:    fmt.Sprintf("%.2f", amount),
		},
	}

	applicationContext := paypal.ApplicationContext{
		ReturnURL: "http://localhost:3000/paypal/confirmpayment",
		CancelURL: "http://localhost:3000/paypal/cancel-payment",
	}
	order, err := client.CreateOrder(
		context.Background(),
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{purchaseUnit},
		nil,
		&applicationContext,
	)
	if err != nil {
		return "", err
	}

	for _, link := range order.Links {
		if link.Rel == "approve" {
			return link.Href, nil
		}
	}

	return "", fmt.Errorf("no approval link found in PayPal response")
}

func CapturePayPalOrder(c *gin.Context) {
	client, err := NewPayPalClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize PayPal client"})
		return
	}

	orderID := c.Query("token")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID missing from query parameters"})
		return
	}
	captureRequest := paypal.CaptureOrderRequest{}

	order, err := client.CaptureOrder(context.Background(), orderID, captureRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to capture PayPal order", "details": err.Error()})
		return
	}

	if order.Status != "COMPLETED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not completed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "order_id": orderID})

}
