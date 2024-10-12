package models

import "time"

type User struct {
	UserID      int `gorm:"primaryKey;autoIncrement"`
	UserName    string
	Address     string
	Email       string `gorm:"not null;unique"`
	Password    string `gorm:"not null"`
	PhoneNumber *string
	Status      string `gorm:"check(status IN('Active', 'Inactive', 'Blocked'))"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	CategoryID   int `gorm:"primaryKey;autoIncrement"`
	CategoryName string
	CreatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}

type Product struct {
	ProductID   int `gorm:"primaryKey;autoIncrement"`
	ProductName string
	Description string
	Price       float64
	CategoryID  int `gorm:"not null;index;constraint:OnDelete:CASCADE;foreignKey:CategoryID;references:CategoryID"`
	ImgURL      string
	Status      string `gorm:"check(status IN('Available', 'Unavailable', 'Out of stock'))"`
	CreatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

type Wishlist struct {
	WishlistID int `gorm:"primaryKey;autoIncrement"`
	ProductID  int `gorm:"not null;foreignKey:ProductID;references:ProductID"`
	UserID     int `gorm:"not null;index;foreignKey:UserID;references:UserID"`
	CreatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}

type Cart struct {
	CartID    int `gorm:"primaryKey;autoIncrement"`
	UserID    int `gorm:"not null;index;foreignKey:UserID;references:UserID"`
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	CartItems []CartItem `gorm:"foreignKey:CartID;references:CartID"`
}

type CartItem struct {
	CartItemID int `gorm:"primaryKey;autoIncrement"`
	CartID     int `gorm:"not null;foreignKey:CartID;references:CartID"`
	ProductID  int `gorm:"not null;foreignKey:ProductID;references:ProductID"`
	Total      int
	Quantity   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Order struct {
	OrderID       int `gorm:"primaryKey;autoIncrement"`
	UserID        int `gorm:"not null;index;foreignKey:UserID;references:UserID"`
	ProductID     int `gorm:"not null;index;foreignKey:ProductID;references:ProductID"`
	PaymentID     int
	OrderDate     time.Time
	Total         int
	CouponID      int `gorm:"foreignKey:CouponID;references:CouponID"`
	Discount      int
	Quantity      int
	Status        string `gorm:"check(status IN('Pending', 'Processing', 'Delivered', 'Canceled'));"`
	Amount        float64
	Method        string `gorm:"check(method IN('Credit Card', 'PayPal', 'Bank Transfer'));"`
	PaymentStatus string `gorm:"check(status IN('Processing', 'Success', 'Failed'));"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID;references:OrderID"`
}

type OrderItem struct {
	OrderItemsID int `gorm:"primaryKey;autoIncrement"`
	OrderID      int `gorm:"not null;index;foreignKey:OrderID;references:OrderID"`
	UserID       int `gorm:"not null;index;foreignKey:UserID;references:UserID"`
	ProductID    int `gorm:"not null;index;foreignKey:ProductID;references:ProductID"`
	Price        int
	Discount     int
}

type Coupon struct {
	CouponID             int    `gorm:"primaryKey;autoIncrement"`
	CouponCode           string `gorm:"unique"`
	CouponDiscountAmount int
	Description          string
	StartDate            time.Time
	Period               int
	MinPurchaseAmount    int
	MaxPurchaseAmount    int
	IsActive             string `gorm:"check(is_active IN('Active', 'Inactive'))"`
}

type ReviewRating struct {
	ReviewRatingID int `gorm:"primaryKey;autoIncrement"`
	UserID         int `gorm:"not null;foreignKey:UserID;references:UserID"`
	ProductID      int `gorm:"not null;foreignKey:ProductID;references:ProductID"`
	Review         string
	Rating         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
