package helper

import (
	"fmt"

	"github.com/go-playground/validator"
)

var validate = validator.New()

func ValidateAll(input any) (string, error) {

	err := validate.Struct(input)
	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "UserName":
				return "username must be alphanumeric and 3-16 characters long", fmt.Errorf("invalid username")
			case "Email":
				return "invalid email format", fmt.Errorf("invalid email")
			case "PhoneNumber":
				return "phone number must be exactly 10 digits", fmt.Errorf("invalid phone number")
			case "Password":
				return "password must be between 8 and 32 characters", fmt.Errorf("invalid password")
			default:
				return "invalid input", fmt.Errorf("validation failed")
			}
		}
	}
	return "", nil
}

func ValidateAddress(input any) (string, error) {

	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "AddressLine1":
				return "street address must be between 3 and 100 characters", fmt.Errorf("invalid street address")
			case "City":
				return "city must be between 2 and 50 characters", fmt.Errorf("invalid city")
			case "PostalCode":
				return "ZIP/postal code must be 5 or 6 characters", fmt.Errorf("invalid ZIP/postal code")
			case "Country":
				return "invalid country code (use 2-letter ISO code)", fmt.Errorf("invalid country code")
			default:
				return "invalid input", fmt.Errorf("validation failed")
			}
		}
	}

	return "", nil
}
