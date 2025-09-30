package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type KeyPool struct {
	db      *sql.DB
	pool    chan string
	minSize int
}

func NewKeyPool(db *sql.DB, size int) *KeyPool {
	log.Printf("Initializing KeyPool with size %d", size)
	keyPool := &KeyPool{
		db:      db,
		pool:    make(chan string, size*2),
		minSize: size / 10,
	}

	go keyPool.refiller()

	return keyPool
}

func (kp *KeyPool) Get() string {
	log.Println("Getting key from pool")
	return <-kp.pool
}

func (kp *KeyPool) Put(key string) {
	log.Printf("Putting key back to pool: %s", key)
	kp.pool <- key
}

func (kp *KeyPool) refiller() {
	log.Println("Starting KeyPool refiller")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for len(kp.pool) < kp.minSize {
			keys, err := kp.fetchKeysFromDB()
			if err != nil {
				log.Printf("Error fetching key from DB: %v", err)
				break
			}

			for _, key := range keys {
				kp.Put(key)
			}
		}
	}
}

func (kp *KeyPool) fetchKeysFromDB() ([]string, error) {
	log.Println("Fetching keys from DB")

	keys := []string{}

	for i := range 100 {
		var key string
		key = "key_" + time.Now().Format("150405.000000000") + "_" + fmt.Sprint(i)
		keys = append(keys, key)
	}

	log.Printf("Fetched %d keys from DB", len(keys))
	return keys, nil
}
