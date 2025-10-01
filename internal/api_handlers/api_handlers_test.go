package api_handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jassi-singh/mini-forge/internal/database"
	"github.com/jassi-singh/mini-forge/internal/repository"
	"github.com/jassi-singh/mini-forge/internal/services"
	"github.com/jassi-singh/mini-forge/internal/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestServer initializes a new server for testing purposes.
func setupTestServer() *httptest.Server {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to in-memory database: " + err.Error())
	}

	database.Migrate(db)

	rangeCounterRepo := repository.NewRangeCounterRepository(db)
	keyPool := services.NewKeyPool(utils.RANGE_SIZE, rangeCounterRepo)
	apiHandler := NewApiHandler(keyPool)

	router := chi.NewRouter()
	router.Get("/get-key", apiHandler.GetKey)

	return httptest.NewServer(router)
}

func TestGetKey_ConcurrentRequests(t *testing.T) {
	// 1. Setup
	server := setupTestServer()
	defer server.Close()

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100, // Allow many idle connections
			MaxIdleConnsPerHost: 100, // For a single host
			IdleConnTimeout:     90 * time.Second,
		},
	}
	// ---------------------------------------------------

	numRequests := 5000 // Let's try 5000 again
	var wg sync.WaitGroup
	wg.Add(numRequests)

	var mu sync.Mutex
	receivedKeys := make(map[string]bool)

	// 2. Act: Simulate concurrent requests
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()

			resp, err := client.Get(server.URL + "/get-key")
			// ----------------------------------------------------
			if err != nil {
				// Use t.Logf for non-fatal errors in goroutines to avoid noisy output
				t.Logf("Failed to make request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK; got %v", resp.Status)
				return
			}

			keyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Logf("Failed to read response body: %v", err)
				return
			}
			key := string(keyBytes)

			mu.Lock()
			if _, exists := receivedKeys[key]; exists {
				t.Errorf("Duplicate key received: %s", key)
			}
			receivedKeys[key] = true
			mu.Unlock()
		}()
	}

	wg.Wait()

	// 3. Assert
	if len(receivedKeys) != numRequests {
		t.Errorf("Expected %d unique keys, but got %d", numRequests, len(receivedKeys))
	}

	t.Logf("Successfully received %d unique keys under concurrent load.", len(receivedKeys))
}

