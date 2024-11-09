package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName    string `gorm:"column:user_name;not null"`
	Email       string `gorm:"column:email;not null"`
	Password    string `gorm:"column:password;not null" json:"-"`
	PhoneNumber string `gorm:"column:phonenumber;not null"`
	Status      string `gorm:"check(status IN('Active', 'Inactive', 'Blocked'))"`
}

type Address struct {
	AddressID    int `gorm:"primaryKey;autoIncrement"`
	UserID       int `gorm:"not null;index;constraint:OnDelete:CASCADE;foreignKey:UserID;references:UserID"`
	AddressLine1 string
	AddressLine2 string
	Country      string
	City         string
	PostalCode   string
	Landmark     string
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Admin struct {
	AdminID   int `gorm:"primaryKey;autoIncrement"`
	AdminName string
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	CategoryID   uint   `gorm:"primaryKey" json:"category_id"`
	CategoryName string `json:"name"`
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Product struct {
	ProductID     int            `gorm:"primaryKey"`
	ProductName   string         `json:"name"`
	Description   string         `json:"description"`
	Price         float64        `json:"price"`
	CategoryID    uint           `gorm:"not null;index;constraint:OnDelete:CASCADE" json:"category_id"`
	ImgURL        string         `json:"img_url"`
	Status        string         `gorm:"check(status IN('Available', 'Out of stock'))" json:"status"`
	Quantity      int            `json:"quantity" gorm:"default:0"`
	AverageRating float64        `gorm:"-" json:"average_rating"`
	TotalReviews  int            `gorm:"-" json:"total_reviews"`
	RecentReviews []ReviewRating `gorm:"foreignKey:ProductID;references:ProductID" json:"recent_reviews"` // Correct reference
	CreatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Wishlist struct {
	WishlistID  int `gorm:"primaryKey;autoIncrement"`
	UserID      int `gorm:"not null;index;foreignKey:UserID;references:UserID"`
	ProductID   int `gorm:"not null;foreignKey:ProductID;references:ProductID"`
	ProductName string
	Price       int
	Quantity    int
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Cart struct {
	CartID    int `gorm:"primaryKey;autoIncrement"`
	UserID    int `gorm:"not null;index"`
	ProductID int `gorm:"not null"`
	Total     int
	Quantity  int
	User      User    `gorm:"foreignKey:UserID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
type Order struct {
	OrderID       int `gorm:"primaryKey;autoIncrement"`
	UserID        int `gorm:"not null;index"`
	PaymentID     string
	OrderDate     time.Time   `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Total         float64     `gorm:"not null"`
	CouponID      int         `gorm:"index"`
	Discount      int         `gorm:"default:0"`
	Quantity      int         `gorm:"default:0"`
	Status        string      `gorm:"check(status IN('Pending', 'Shipped', 'Delivered', 'Canceled','Failed','Returned'))"`
	Method        string      `gorm:"check(method IN('Credit Card', 'PayPal', 'Bank Transfer'))"`
	PaymentStatus string      `gorm:"check(payment_status IN('Processing', 'Success', 'Failed'))"`
	CreatedAt     time.Time   `gorm:"autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime"`
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	OrderItemsID int     `gorm:"primaryKey;autoIncrement"`
	OrderID      int     `gorm:"not null;index"`
	UserID       int     `gorm:"not null;index"`
	ProductID    int     `gorm:"not null;index"`
	Quantity     int     `gorm:"default:0"`
	Price        float64 `gorm:"not null"`
	Discount     int     `gorm:"default:0"`
}

type Coupon struct {
	CouponID          int    `gorm:"primaryKey;autoIncrement"`
	CouponCode        string `gorm:"unique" json:"coupon_code"`
	DiscountAmount    float64
	DiscountType      string `gorm:"not null;default:fixed"`
	Description       string
	StartDate         time.Time
	EndDate           time.Time
	MinPurchaseAmount int
	MaxPurchaseAmount int
	IsActive          bool
	DeletedAt         gorm.DeletedAt
}
type ReviewRating struct {
	ReviewRatingID int       `gorm:"primaryKey;autoIncrement" json:"review_rating_id"`
	UserID         int       `gorm:"not null;index" json:"user_id"`    // Foreign key reference to User
	ProductID      int       `gorm:"not null;index" json:"product_id"` // Foreign key reference to Product
	Rating         int       `json:"rating"`
	Comment        string    `json:"comment"` // Corrected JSON tag
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID;references:ID"`           // Assuming ID is the primary key in User
	Product Product `gorm:"foreignKey:ProductID;references:ProductID"` // Assuming ProductID is the primary key in Product
}
type TempUser struct {
	UserName    string `json:"username"`
	Address     string
	Email       string `json:"email"`
	Password    string
	PhoneNumber string
}

type UserLoginMethod struct {
	UserLoginMethodEmail string
	LoginMethod          string
}

type OTP struct {
	Email  string
	Code   string
	Expiry time.Time
}

type Wallet struct {
	WalletID  uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null;foreignKey:UserID;references:UserID"`
	Balance   float64 `gorm:"not null;default:0.0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Offer struct {
	gorm.Model
	ProductID       int `gorm:"not null"`
	OfferPercentage int `gorm:"not null"`
}

type SalesReport struct {
	DateRange        string  `json:"date_range"`
	TotalSalesCount  int     `json:"total_sales_count"`
	TotalOrderAmount float64 `json:"total_order_amount"`
	TotalDiscount    float64 `json:"total_discount"`
	CouponsDeduction float64 `json:"coupons_deduction"`
}
type TempOrder struct {
	OrderID       string    `gorm:"primaryKey;type:varchar(255);not null"`
	UserID        int       `gorm:"column:user_id"`
	PaymentID     int       `gorm:"column:payment_id"`
	OrderDate     time.Time `gorm:"column:order_date"`
	Total         float64   `gorm:"column:total"`
	CouponID      int       `gorm:"column:coupon_id"`
	Discount      int       `gorm:"column:discount"`
	Quantity      int       `gorm:"column:quantity"`
	Status        string    `gorm:"column:status"`
	Method        string    `gorm:"column:method"`
	PaymentStatus string    `gorm:"column:payment_status"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

type Invoice struct {
	InvoiceID      string        `json:"invoice_id"`
	Date           time.Time     `json:"date"`
	UserID         int           `json:"user_id"`
	BillingAddress string        `json:"billing_address"`
	Items          []InvoiceItem `json:"items"`
	Subtotal       float64       `json:"subtotal"`
	Tax            float64       `json:"tax"`
	Discount       int           `json:"discount"`
	Total          float64       `json:"total"`
}

type InvoiceItem struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

type WalletTransaction struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null"`
	OrderID         uint      `gorm:"default:null"`
	Amount          float64   `gorm:"not null"`
	TransactionType string    `gorm:"type:varchar(10);not null"`
	Description     string    `gorm:"type:varchar(255)"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}
