package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jassi-singh/mini-forge/internal/api_handlers"
	"github.com/jassi-singh/mini-forge/internal/database"
	"github.com/jassi-singh/mini-forge/internal/keypool"
)

func main() {
	log.Print("Initializing application...")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	database.Migrate(db)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	keyPool := utils.NewKeyPool(100)

	apiHandler := api_handlers.NewApiHandler(keyPool)

	router.Get("/get-key", apiHandler.GetKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	// Start the HTTP server

	log.Println("Starting HTTP server on :", port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
