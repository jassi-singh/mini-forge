# mini-forge

A lightweight, high-performance, and scalable Key Generation Service (KGS) written in Go. Designed for distributed systems that require a steady stream of short, unique, and collision-free keys.

## ✨ Key Features

*   High-Performance: Serves pre-generated keys directly from an in-memory pool, eliminating database contention on read operations.
    
*   Scalable & Distributed-Ready: Multiple instances can run concurrently without generating duplicate keys, thanks to a centralized range-based counter.
    
*   Collision-Free by Design: Guarantees unique keys across all instances by using a transactional database mechanism to reserve blocks of numbers for key generation.
    
*   Concurrent-Safe: Built to handle thousands of simultaneous requests safely using a channel-based key pool and atomic database operations.
    

## ⚙️ How It Works

The service is designed to avoid the typical bottlenecks of a traditional KGS. Instead of hitting the database for every key request, it generates keys in batches and serves them from memory.

1.  Range Reservation: When a mini-forge instance needs new keys, it requests a unique range of numbers from the database. It does this in a transaction using a `SELECT ... FOR UPDATE` lock to ensure no two instances can get the same range. The size of this range is configurable via the RANGE\_SIZE setting.
    
2.  Key Generation: The service then generates keys for every number in its reserved range using a Base62 encoding strategy. This creates short, URL-friendly strings.
    
3.  In-Memory Pool: These generated keys are stored in an in-memory channel that acts as a key pool.
    
4.  Serving Keys: When a client requests a key, the service simply takes one from the in-memory pool—a very fast operation.
    
5.  Automatic Refill: A background goroutine continuously monitors the size of the key pool. When the number of available keys drops below a certain threshold, it automatically fetches a new range from the database and refills the pool, ensuring there's always a ready supply of keys.
    

This architecture minimizes database interaction, making the service extremely fast and scalable.

## 🚀 Getting Started

### Prerequisites

*   Go 1.19 or higher
    
*   SQLite3
    

### Configuration

The service is configured using a config.yml file, with values populated by environment variables.

1.  Create a .env file in the root directory by copying the example:  
    cp .env.example .env  
      
    
2.  Edit the .env file with your configuration:  
    ```
    # The port the HTTP server will listen on  
    PORT=8080  
      
    # The number of keys to pre-generate and pool in memory at a time.  
    # A larger range size means fewer database queries but higher initial memory use.  
    RANGE_SIZE=1000  
      
    # Set to true to enable detailed debug logging  
    DEBUG_ENABLED=false  
      
    # Absolute or relative path to the configuration file  
    CONFIG_PATH=./config/config.yml 
    ``` 
      
    

### Build and Run

1.  Clone the repository:  
    ```
        git clone https://github.com/jassi-singh/mini-forge.git
        cd mini-forge  
    ```
      
    
2.  Install dependencies:  
    ```
    go mod download  
    ```
      
    
3.  Run the service:  
    ```
    go run main.go  
      
    #Or build a binary and run it:  
    go build -o mini-forge  
    ./mini-forge  
    ```
    

The service will start on the port configured in your .env file.

## API Usage

The service exposes a single, simple endpoint for retrieving a key.

### GET /get-key

Returns a unique, Base62-encoded key.

Example Request:
```bash
curl http://localhost:8080/get-key 
``` 
  

Success Response (200 OK):
```
000001a  
```
  

## 🧪 Testing the Service

The project includes a comprehensive concurrency test to simulate a high-traffic environment and ensure no duplicate keys are ever generated.

### ⚠️ Important: System Limits for Concurrent Testing

When running tests with thousands of concurrent workers, you may hit your system's file descriptor limit. This can cause "too many open files" errors.

**Check your current limit:**
```bash
ulimit -n
```

**Temporarily increase the limit (for the current session):**
```bash
# Increase to 10,000 file descriptors
ulimit -n 10000
```

**For macOS users, you may need to check both soft and hard limits:**
```bash
# Check soft limit
ulimit -Sn

# Check hard limit
ulimit -Hn

# Set to the hard limit or a specific value
ulimit -Sn 10000
```

> **Note:** The required limit depends on the number of concurrent workers in your test. For very high concurrency (10,000+ workers), you may need to increase this further.

#### Running the Test

```
# Run all tests, including the concurrency test  
go test -v ./...  
  
# Run with a longer timeout if needed  
go test -v -timeout 30s ./internal/api\_handlers/  
```
  

The test simulates thousands of concurrent requests and verifies that every single key received is unique. This is crucial for validating the service's reliability for production use.