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
	err = Db.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Admin{},
		&models.Category{},  
		&models.Product{}, 
		&models.Wishlist{},
		&models.Order{},
		&models.OrderItem{},
		&models.Coupon{},
		&models.ReviewRating{},
		&models.OTP{},
		&models.TempUser{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

}
