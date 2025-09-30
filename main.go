package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jassi-singh/mini-forge/internal/api_handlers"
	"github.com/jassi-singh/mini-forge/internal/keypool"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Print("Initializing application...")

	db, err := sql.Open("sqlite3", "../app.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	keyPool := utils.NewKeyPool(db, 100)

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
