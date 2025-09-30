package database

import (
	"log"

	"github.com/jassi-singh/mini-forge/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
}

func Migrate(db *gorm.DB) {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(&models.RangeCounter{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	initialCounter := models.RangeCounter{ID: 1, LastUsed: 0}
	result := db.FirstOrCreate(&initialCounter, models.RangeCounter{ID: 1})
	if result.Error != nil {
		log.Fatalf("Failed to initialize RangeCounter: %v", result.Error)
	}

	log.Println("Database migrations completed.")
}
