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

### Basic Test

Once the service is running, you can test it using `curl`:

```bash
# Get a unique key
curl http://localhost:8080/get-key
```

**Expected Response:**
```
aB3xY9
```

### Load Testing

Test the service under load using multiple concurrent requests:

```bash
# Using ab (Apache Bench)
ab -n 10000 -c 100 http://localhost:8080/get-key

# Using curl in a loop
for i in {1..10}; do
  curl http://localhost:8080/get-key
  echo ""
done
```

### Verify Uniqueness

Generate multiple keys and verify they are unique:

```bash
# Generate 100 keys and check for duplicates
for i in {1..100}; do
  curl -s http://localhost:8080/get-key
  echo ""
done | sort | uniq -d
```

If the output is empty, all keys are unique! 🎉

## 📝 Logging

The application uses a custom logger with three log levels:

* **INFO** - Always logged, used for important application events
* **ERROR** - Always logged, used for error conditions
* **DEBUG** - Only logged when `DEBUG_ENABLED=true` is set in the environment

To enable debug logging, set the `DEBUG_ENABLED` environment variable:
```bash
export DEBUG_ENABLED=true
```


