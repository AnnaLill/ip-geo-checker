package repository

import (
	"ip-geo-checker/internal/models"
	"sync"
)

// InMemoryCache реализация CacheRepository с использованием in-memory хранилища
type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]models.CheckIPResponse
}

// NewInMemoryCache создает новый экземпляр InMemoryCache
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]models.CheckIPResponse),
	}
}

// Get получает значение из кэша
func (c *InMemoryCache) Get(ip string) (models.CheckIPResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.items[ip]
	return value, exists
}

// Set сохраняет значение в кэш
func (c *InMemoryCache) Set(ip string, response models.CheckIPResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[ip] = response
}

// Count возвращает количество элементов в кэше
func (c *InMemoryCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
