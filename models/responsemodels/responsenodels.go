package responsemodels

import (
	"time"
)

type User struct {
	UserName    string `gorm:"column:user_name;not null"`
	Email       string `gorm:"column:email;not null"`
	PhoneNumber string `gorm:"column:phonenumber;not null"`
	Status      string `gorm:"check(status IN('Active', 'Inactive', 'Blocked'))"`
}

type Products struct {
	ProductID   int     `gorm:"primaryKey"`
	ProductName string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  uint    `gorm:"not null;index;constraint:OnDelete:CASCADE" json:"category_id"`
	ImgURL      string  `json:"img_url"`
	Status      string  `gorm:"check(status IN('Available', 'Out of stock'))"`
	Quantity    int     `json:"quantity"`
}
type Address struct {
	AddressID    int    `gorm:"primaryKey;autoIncrement"`
	AddressLine1 string `json:"addressline1"`
	AddressLine2 string `json:"addressline2"`
	Country      string `json:"country"`
	City         string `json:"city"`
	PostalCode   string `json:"postalcode"`
	Landmark     string `json:"landmark"`
}
type CartResponse struct {
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Total     int    `json:"total"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
}

type OrderResponse struct {
	UserID         int       `json:"user_id"`
	OrderID        int       `json:"order_id"`
	Total          float64   `json:"total"`
	Quantity       int       `json:"quantity"`
	DiscountAmount float64   `json:"discount_amount"`
	Status         string    `json:"status"`
	Method         string    `json:"method"`
	PaymentStatus  string    `json:"payment_status"`
	OrderDate      time.Time `json:"order_date"`
}

type Wishlist struct {
	UserID      int    `json:"user_id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
}
