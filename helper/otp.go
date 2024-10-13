package helper

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
)

func GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			log.Fatal("Error Generating OTP", err)
		}
		otp += fmt.Sprintf("%d", digit)
	}
	return otp, nil
}

func SendEmail(to, otp string) {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	if from =="" || password
}
