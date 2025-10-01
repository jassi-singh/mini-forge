package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jassi-singh/mini-forge/internal/api_handlers"
	"github.com/jassi-singh/mini-forge/internal/database"
	"github.com/jassi-singh/mini-forge/internal/repository"
	"github.com/jassi-singh/mini-forge/internal/services"
	"github.com/jassi-singh/mini-forge/internal/utils"
)

func main() {
	log.Print("Initializing application...")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	database.Migrate(db)

	config, err := utils.LoadConfig("./config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	rangeCounterRepo := repository.NewRangeCounterRepository(db, config)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	keyPool := services.NewKeyPool(config.RangeSize, rangeCounterRepo, config)

	apiHandler := api_handlers.NewApiHandler(keyPool)

	router.Get("/get-key", apiHandler.GetKey)

	addr := ":" + strconv.Itoa(config.Port)
	// Start the HTTP server

	log.Println("Starting HTTP server on :", config.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
