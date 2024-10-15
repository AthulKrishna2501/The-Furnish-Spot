package models

type SignupInput struct {
	UserName    string `json:"username" validate:"required,min=3,max=16,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phonenumber" validate:"required,len=10,numeric"`
	Password    string `json:"password" validate:"required,min=8,max=32"`
}

type VerifyOTP struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
