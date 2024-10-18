package responsemodels

import "time"

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
