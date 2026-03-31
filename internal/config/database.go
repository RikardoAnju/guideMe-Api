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
		&models.Destinasi{},
		&models.Event{},
	)

	// Tambah index untuk performa query users
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("failed to get sql.DB:", err)
		return
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`,
		`CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
	}

	for _, q := range indexes {
		if _, err := sqlDB.Exec(q); err != nil {
			log.Printf("index warning: %v", err)
		}
	}

	log.Println("Migrations completed")
}