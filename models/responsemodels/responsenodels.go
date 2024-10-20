package responsemodels

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ProductID   int     `gorm:"primaryKey"`
	ProductName string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  uint    `gorm:"not null;index;constraint:OnDelete:CASCADE;foreignKey:CategoryID;references:CategoryID" json:"category_id"`
	ImgURL      string  `json:"img_url"`
	Status      string  `gorm:"check(status IN('Available', 'Unavailable', 'Out of stock'))"`
	CreatedAt   time.Time
}

type User struct {
	gorm.Model
	UserName    string `gorm:"column:user_name;not null"`
	Email       string `gorm:"column:email;not null"`
	PhoneNumber string `gorm:"column:phonenumber;not null"`
	Status      string `gorm:"check(status IN('Active', 'Inactive', 'Blocked'))"`
}
