package config

import (
	"log"
	"os"

	"guide-me/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func RunMigrations() {
	DB.AutoMigrate(
		&models.User{},
		&models.TokenBlacklist{},
	)
	log.Println("Migrations completed")
}
