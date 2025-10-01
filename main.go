package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jassi-singh/mini-forge/internal/api_handlers"
	"github.com/jassi-singh/mini-forge/internal/database"
	"github.com/jassi-singh/mini-forge/internal/logger"
	"github.com/jassi-singh/mini-forge/internal/repository"
	"github.com/jassi-singh/mini-forge/internal/services"
	"github.com/jassi-singh/mini-forge/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		logger.Fatal("CONFIG_PATH environment variable is not set")
	}
	// Initialize logger with default settings first
	logger.InitLogger(false)

	logger.Info("Initializing application...")

	db, err := database.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	database.Migrate(db)

	config, err := utils.LoadConfig(configPath)
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	// Reinitialize logger with debug setting from config
	logger.InitLogger(config.DebugEnabled)
	logger.Debug("Debug logging is enabled")

	rangeCounterRepo := repository.NewRangeCounterRepository(db, config)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	keyPool := services.NewKeyPool(config.RangeSize, rangeCounterRepo, config)

	apiHandler := api_handlers.NewApiHandler(keyPool)

	router.Get("/get-key", apiHandler.GetKey)

	addr := ":" + strconv.Itoa(config.Port)
	// Start the HTTP server

	logger.Info("Starting HTTP server on :%d", config.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("Failed to start HTTP server: %v", err)
	}
}
