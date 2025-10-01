package database

import (
	"github.com/jassi-singh/mini-forge/internal/logger"
	"github.com/jassi-singh/mini-forge/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
}

func Migrate(db *gorm.DB) {
	logger.Info("Running database migrations...")

	err := db.AutoMigrate(&models.RangeCounter{})
	if err != nil {
		logger.Fatal("Failed to migrate database: %v", err)
	}

	initialCounter := models.RangeCounter{ID: 1, LastUsed: 0}
	result := db.FirstOrCreate(&initialCounter, models.RangeCounter{ID: 1})
	if result.Error != nil {
		logger.Fatal("Failed to initialize RangeCounter: %v", result.Error)
	}

	logger.Info("Database migrations completed.")
}
