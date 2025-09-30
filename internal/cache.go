package internal

type CacheStore interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

type InMemoryCache struct {
	store map[string]string
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		store: make(map[string]string),
	}
}

func (c *InMemoryCache) Get(key string) (string, bool) {
	value, exists := c.store[key]
	return value, exists
}

func (c *InMemoryCache) Set(key string, value string) {
	c.store[key] = value
}

func (c *InMemoryCache) Delete(key string) {
	delete(c.store, key)
}
