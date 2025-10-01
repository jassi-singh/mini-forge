package services

import (
	"sync"
	"time"

	"github.com/jassi-singh/mini-forge/internal/logger"
	"github.com/jassi-singh/mini-forge/internal/repository"
	"github.com/jassi-singh/mini-forge/internal/utils"
)

type KeyPool struct {
	config           *utils.Config
	rangeCounterRepo repository.RangeCounterRepository
	pool             chan string
	minSize          int
	refillerMutex    sync.Mutex
}

func NewKeyPool(size int, rangeCounterRepo repository.RangeCounterRepository, config *utils.Config) *KeyPool {
	logger.Info("Initializing KeyPool with size %d", size)
	keyPool := &KeyPool{
		config:           config,
		rangeCounterRepo: rangeCounterRepo,
		pool:             make(chan string, size*2),
		minSize:          size / 10,
		refillerMutex:    sync.Mutex{},
	}

	go keyPool.refiller()

	return keyPool
}

func (kp *KeyPool) Get() string {
	logger.Debug("Getting key from pool")
	return <-kp.pool
}

func (kp *KeyPool) Put(key string) {
	logger.Debug("Putting key back to pool: %s", key)
	kp.pool <- key
}

func (kp *KeyPool) refiller() {
	logger.Debug("Starting KeyPool refiller")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if len(kp.pool) <= kp.minSize {
			kp.refillerMutex.Lock()
			if len(kp.pool) <= kp.minSize {

				keys, err := kp.fetchKeysFromDB()
				if err != nil {
					logger.Error("Error fetching key from DB: %v", err)
					kp.refillerMutex.Unlock()
					continue
				}

				kp.refillerMutex.Unlock()
				for _, key := range keys {
					kp.Put(key)
				}
			} else {
				kp.refillerMutex.Unlock()
			}
		}
	}
}

func (kp *KeyPool) fetchKeysFromDB() ([]string, error) {
	logger.Debug("Fetching keys from DB")

	keys := []string{}

	lastUsed, err := kp.rangeCounterRepo.GetAndIncrement()
	if err != nil {
		return nil, err
	}

	for i := lastUsed; i < lastUsed+int64(kp.config.RangeSize); i++ {
		key := utils.GenerateBase62Key(i)
		keys = append(keys, key)
	}

	logger.Debug("Fetched %d keys from DB", len(keys))
	return keys, nil
}
