package db

import (
	"log"
	"os"

	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env", err)
	}
	Db, err = gorm.Open(postgres.Open(os.Getenv("dsn")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error loading database", err)
		return
	}
	Db.AutoMigrate(&models.User{}, &models.Address{}, &models.Admin{}, &models.Category{}, &models.Product{}, &models.Wishlist{}, &models.Cart{}, &models.CartItem{}, &models.Order{}, &models.OrderItem{}, &models.Coupon{}, &models.ReviewRating{}, models.OTP{}, models.TempUser{})

}
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Admin{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Coupon{})
	db.AutoMigrate(&models.Address{})
	db.AutoMigrate(&models.Wishlist{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.CartItem{})
	db.AutoMigrate(&models.Order{})
	db.AutoMigrate(&models.OrderItem{})
	db.AutoMigrate(&models.ReviewRating{})
	db.AutoMigrate(&models.OTP{})
}
