package models

import "time"

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
	ReEnter     string `json:"reenter" validate:"required,min=8,max=32"`
}

type InputAddress struct {
	AddressLine1 string `json:"addressline1"`
	AddressLine2 string `json:"addressline2"`
	Country      string `json:"country" validate:"required,oneof=IND US UK CA AU"`
	City       string `json:"city"`
	PostalCode string `json:"postalcode" validate:"required,len=6,numeric"`
	Landmark   string `json:"landmark"`
}

type OrderInput struct {
	AddressID  int    `json:"address_id"`
	CouponCode string `json:"coupon_code,omitempty"`
	Method     string `json:"method"`
}

type CartInput struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
type ReviewInput struct {
	ProductID int    `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

type CouponInput struct {
	CouponCode        string    `json:"coupon_code"`
	DiscountAmount    float64   `json:"discount_amount"`
	DiscountType      string    `json:"discount_type"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	MinPurchaseAmount int       `json:"min_purchase_amount"`
	MaxPurchaseAmount int       `json:"max_purchase_amount"`
	IsActive          bool      `json:"is_active"`
}

type ReturnOrder struct {
	OrderID int    `json:"order_id"`
	Reason  string `json:"reason"`
}

type OfferInput struct {
	ProductID       int `json:"product_id"`
	OfferPercentage int `json:"offer_percentage"`
}
