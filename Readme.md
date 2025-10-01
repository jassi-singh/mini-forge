# mini-forge

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/mini-forge)](https://goreportcard.com/report/github.com/your-username/mini-forge)
A lightweight, high-performance Key Generation Service written in Go, designed for distributed systems.

## ✨ Key Features

* **High-Performance:** Serves pre-generated keys from memory with zero database contention on reads.
* **Scalable:** Designed to run with multiple instances without conflict.
* **Collision-Free:** Uses a range-based, Base-62 encoding strategy to guarantee unique keys.
* **Concurrent-Safe:** Built to handle thousands of simultaneous requests.

## 📝 Logging

The application uses a custom logger with three log levels:

* **INFO** - Always logged, used for important application events
* **ERROR** - Always logged, used for error conditions
* **DEBUG** - Only logged when `DEBUG_ENABLED=true` is set in the environment

To enable debug logging, set the `DEBUG_ENABLED` environment variable:
```bash
export DEBUG_ENABLED=true
```


