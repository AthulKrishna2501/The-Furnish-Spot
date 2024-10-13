package helper

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
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

	if from == "" || password == "" {
		log.Fatal("Email credentials not set in environment variable")
	}

	msg := "Your OTP for Signup is" + otp

	err := smtp.SendMail("smtp.example.com:587", smtp.PlainAuth("", from, password, "smtp.example.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Fatal("Error Sending otp", err)
	}
	fmt.Println("OTP sent successfully")
}
