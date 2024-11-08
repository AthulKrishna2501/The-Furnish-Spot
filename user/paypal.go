package user

import (
	"context"
	"fmt"
	"log"
	"os"

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

func CreatePayPalPayment(client *paypal.Client, amount float64) (string, string, error) {
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
		return "", "", err
	}

	var approvalURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}

	if approvalURL == "" {
		return "", "", fmt.Errorf("no approval link found in PayPal response")
	}

	return approvalURL, order.ID, nil
}
