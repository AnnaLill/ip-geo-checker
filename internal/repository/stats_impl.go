package repository

import (
	"ip-geo-checker/internal/models"
	"sync"
)

// InMemoryStats реализация StatsRepository с использованием in-memory хранилища
type InMemoryStats struct {
	mu                              sync.RWMutex
	requestsCount                   int
	cacheCount                      int
	externalSuccessfulRequestsCount int
}

// NewInMemoryStats создает новый экземпляр InMemoryStats
func NewInMemoryStats() *InMemoryStats {
	return &InMemoryStats{}
}

// IncrementRequests увеличивает счетчик запросов
func (s *InMemoryStats) IncrementRequests() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestsCount++
}

// IncrementCache увеличивает счетчик попаданий в кэш
func (s *InMemoryStats) IncrementCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cacheCount++
}

// IncrementExternalRequests увеличивает счетчик успешных внешних запросов
func (s *InMemoryStats) IncrementExternalRequests(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.externalSuccessfulRequestsCount += count
}

// GetStats возвращает текущую статистику
func (s *InMemoryStats) GetStats() models.StatsResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return models.StatsResponse{
		RequestsCount:                   s.requestsCount,
		CacheCount:                      0, // Перезаписывается в service
		ExternalSuccessfulRequestsCount: s.externalSuccessfulRequestsCount,
	}
}
