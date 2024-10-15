package helper

import (
	"fmt"

	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/go-playground/validator"
)

var validate = validator.New()

func ValidateAll(input models.SignupInput) (string, error) {

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
