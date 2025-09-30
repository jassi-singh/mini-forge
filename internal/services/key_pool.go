package services

import (
	"log"
	"time"

	"github.com/jassi-singh/mini-forge/internal/repository"
	"github.com/jassi-singh/mini-forge/internal/utils"
)

type KeyPool struct {
	rangeCounterRepo repository.RangeCounterRepository
	pool             chan string
	minSize          int
}

func NewKeyPool(size int, rangeCounterRepo repository.RangeCounterRepository) *KeyPool {
	log.Printf("Initializing KeyPool with size %d", size)
	keyPool := &KeyPool{
		rangeCounterRepo: rangeCounterRepo,
		pool:             make(chan string, size*2),
		minSize:          size / 10,
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

	counter, err := kp.rangeCounterRepo.GetAndIncrement()
	if err != nil {
		return nil, err
	}

	for i := range counter {
		key := utils.GenerateBase62Key(i)
		keys = append(keys, key)
	}

	log.Printf("Fetched %d keys from DB", len(keys))
	return keys, nil
}
