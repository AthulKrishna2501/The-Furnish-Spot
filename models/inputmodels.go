package models

type SignupInput struct {
	UserName    string `json:"username" validate:"required,min=3,max=16,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phonenumber" validate:"required,len=10,numeric"`
	Password    string `json:"password" validate:"required,min=8,max=32"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	Captcha     string `json:"captcha" binding:"required"`
}

type VerifyOTP struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type SearchProduct struct {
	Name string `json:"name" binding:"required"`
}

type EditUser struct {
	UserName    string `json:"username" validate:"required,min=3,max=16,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phonenumber" validate:"required,len=10,numeric"`
	Password    string `json:"password" validate:"required,min=8,max=32"`
}

type NewPassword struct {
	Password    string `json:"password" validate:"required,min=8,max=32"`
	NewPassword string `json:"newpassword" validate:"required,min=8,max=32"`
}

type InputAddress struct {
	AddressLine1 string `json:"addressline1"`
	AddressLine2 string `json:"addressline2"`
	Country      string `json:"country"`
	City         string `json:"city"`
	PostalCode   string `json:"postalcode"`
	Landmark     string `json:"landmark"`
}

type OrderInput struct {
	AddressID int `json:"address_id"`
}
