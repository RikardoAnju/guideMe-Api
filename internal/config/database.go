package config

import (
	"guide-me/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dbURL := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	DB = db
	log.Println("Connected to Supabase PostgreSQL 🚀")
}

func RunMigrations() {

	err := DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrations applied successfully ✅")
}
