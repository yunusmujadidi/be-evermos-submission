package database

import (
	"be-evermos-submission/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	DB.AutoMigrate(
		&models.User{},
		&models.Address{}, 
		&models.Store{},
		&models.Category{},
		&models.Product{},
		&models.ProductPhoto{},
		&models.Transaction{},
		&models.TransactionDetail{},
		&models.ProductLog{},
	)
	
	log.Println("Database connected and migrated!")
} 