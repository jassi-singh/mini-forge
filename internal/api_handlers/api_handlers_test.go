package api_handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jassi-singh/mini-forge/internal/api_handlers"
	"github.com/jassi-singh/mini-forge/internal/database"
	"github.com/jassi-singh/mini-forge/internal/repository"
	"github.com/jassi-singh/mini-forge/internal/services"
	"github.com/jassi-singh/mini-forge/internal/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestServer initializes a new server for testing purposes.
func setupTestServer() *httptest.Server {
	// Load environment variables from .env file
	_ = godotenv.Load("../../.env")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to in-memory database: " + err.Error())
	}

	database.Migrate(db)
	config, err := utils.LoadConfig("../../config/config.yml")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	rangeCounterRepo := repository.NewRangeCounterRepository(db, config)
	keyPool := services.NewKeyPool(config.RangeSize, rangeCounterRepo, config)
	apiHandler := api_handlers.NewApiHandler(keyPool)

	router := chi.NewRouter()
	router.Get("/get-key", apiHandler.GetKey)

	return httptest.NewServer(router)
}

func TestGetKey_ConcurrencyWithSyncMap(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	// Create a test server

	// Number of concurrent requests
	numRequests := 5000

	// Use a sync.Map to store generated keys
	var keys sync.Map

	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()

			resp, err := http.Get(server.URL + "/get-key")
			if err != nil {
				t.Errorf("Request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %v", err)
				return
			}
			key := string(body)

			// LoadOrStore is an atomic operation that checks for the key's
			// existence and stores it if it's not there.
			// It returns the existing value if the key is already present.
			_, loaded := keys.LoadOrStore(key, true)
			if loaded {
				t.Errorf("Duplicate key generated: %s", key)
			}
		}()
	}

	wg.Wait()
}
