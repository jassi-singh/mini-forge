package internal

import "sync"

type CacheStore interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

type InMemoryCache struct {
	store sync.Map
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{}
}

func (c *InMemoryCache) Get(key string) (string, bool) {
	if value, exists := c.store.Load(key); exists {
		return value.(string), true
	}
	return "", false
}

func (c *InMemoryCache) Set(key string, value string) {
	c.store.Store(key, value)
}

func (c *InMemoryCache) Delete(key string) {
	c.store.Delete(key)
}
