package api_handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
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

	numRequests := 5000
	maxConcurrentWorkers := 1000
	var keys sync.Map
	var failedRequests atomic.Int64
	var duplicateKeys atomic.Int64

	// Create a channel for jobs (the requests to be made)
	jobs := make(chan int, numRequests)
	for i := range numRequests {
		jobs <- i
	}
	close(jobs)

	var wg sync.WaitGroup

	// Start a fixed number of worker goroutines
	for range maxConcurrentWorkers {
		wg.Go(func() {
			// Each worker pulls jobs from the channel until it's empty
			for range jobs {
				resp, err := http.Get(server.URL + "/get-key")
				if err != nil {
					failedRequests.Add(1)
					continue
				}

				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					failedRequests.Add(1)
					continue
				}
				key := string(body)

				_, loaded := keys.LoadOrStore(key, true)
				if loaded {
					duplicateKeys.Add(1)
				}
			}
		})
	}

	wg.Wait()

	// Report the results
	failedCount := failedRequests.Load()
	duplicateCount := duplicateKeys.Load()
	successCount := numRequests - int(failedCount)

	t.Logf("Total requests: %d", numRequests)
	t.Logf("Successful requests: %d", successCount)
	t.Logf("Failed requests: %d", failedCount)
	t.Logf("Duplicate keys: %d", duplicateCount)

	// Assertions
	if duplicateCount > 0 {
		t.Errorf("Expected no duplicate keys, but found %d duplicates", duplicateCount)
	}
}
