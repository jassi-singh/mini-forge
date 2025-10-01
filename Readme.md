# mini-forge

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/mini-forge)](https://goreportcard.com/report/github.com/your-username/mini-forge)
A lightweight, high-performance Key Generation Service written in Go, designed for distributed systems.

## ✨ Key Features

* **High-Performance:** Serves pre-generated keys from memory with zero database contention on reads.
* **Scalable:** Designed to run with multiple instances without conflict.
* **Collision-Free:** Uses a range-based, Base-62 encoding strategy to guarantee unique keys.
* **Concurrent-Safe:** Built to handle thousands of simultaneous requests.


## 🚀 Running the Service

### Prerequisites

* Go 1.19 or higher
* SQLite (included with most systems)

### Configuration

The service requires the following environment variables:

```bash
export CONFIG_PATH=./config/config.yml
export PORT=8080
export RANGE_SIZE=1000
export DEBUG_ENABLED=false  # Optional: set to true for debug logs
```

### Build and Run

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jassi-singh/mini-forge.git
   cd mini-forge
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set environment variables:**
   ```bash
   export CONFIG_PATH=./config/config.yml
   export PORT=8080
   export RANGE_SIZE=1000
   export DEBUG_ENABLED=false
   ```

4. **Run the service:**
   ```bash
   go run main.go
   ```

   Or build and run:
   ```bash
   go build -o mini-forge
   ./mini-forge
   ```

The service will start on the configured port (default: 8080).

## 🧪 Testing the Service

### Concurrency Test

The service includes a comprehensive concurrency test that simulates real-world high-traffic scenarios to verify that:
- ✅ Keys are generated correctly under heavy load
- ✅ No duplicate keys are ever produced
- ✅ The system handles thousands of concurrent requests safely

#### Running the Test

```bash
# Run the concurrency test
go test -v ./internal/api_handlers/

# Run with timeout (recommended for large tests)
go test -v -timeout 30s ./internal/api_handlers/
```

#### Understanding the Test Configuration

The test is configured with two main parameters:

**1. `numRequests` (default: 5000)**
- **What it means:** Total number of key generation requests to make
- **Web server analogy:** Total number of HTTP requests your API will receive during a traffic spike
- **Example:** If set to 5000, the test simulates 5000 clients requesting unique keys

**2. `maxConcurrentWorkers` (default: 1000)**
- **What it means:** Number of goroutines (workers) making requests simultaneously
- **Web server analogy:** Maximum number of active connections your server handles at once. Even if 5000 requests come in, only 1000 are processed concurrently—similar to how web servers limit concurrent connections to manage resources
- **Example:** With 1000 workers and 5000 requests, workers continuously pull new requests from the queue until all 5000 are complete

#### Test Output

After running, you'll see a summary like this:

```
=== RUN   TestGetKey_ConcurrencyWithSyncMap
    api_handlers_test.go:101: Total requests: 5000
    api_handlers_test.go:102: Successful requests: 5000
    api_handlers_test.go:103: Failed requests: 0
    api_handlers_test.go:104: Duplicate keys: 0
--- PASS: TestGetKey_ConcurrencyWithSyncMap (2.34s)
PASS
```

**What each metric means:**
- **Total requests:** How many key generation requests were attempted
- **Successful requests:** Requests that completed successfully
- **Failed requests:** Requests that failed (network errors, timeouts, etc.)
- **Duplicate keys:** Number of duplicate keys found (should ALWAYS be 0)

#### How It Works (Real-World Scenario)

Think of it like a web API under heavy load:

1. **5000 incoming requests** hit your API endpoint (`numRequests`)
2. Your server accepts up to **1000 concurrent connections** at a time (`maxConcurrentWorkers`)
3. All connections are processed in parallel (true concurrency)
4. Each request receives a unique key in the response
5. After all requests complete, we verify that no two responses contained the same key

If even one duplicate is found, the test fails! This ensures your key generation service is **production-ready** and can handle real-world traffic without collisions. 🎯

## 📝 Logging

The application uses a custom logger with three log levels:

* **INFO** - Always logged, used for important application events
* **ERROR** - Always logged, used for error conditions
* **DEBUG** - Only logged when `DEBUG_ENABLED=true` is set in the environment

To enable debug logging, set the `DEBUG_ENABLED` environment variable:
```bash
export DEBUG_ENABLED=true
```


