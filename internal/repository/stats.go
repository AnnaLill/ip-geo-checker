package repository

import "ip-geo-checker/internal/models"

// StatsRepository интерфейс для работы со статистикой
type StatsRepository interface {
	IncrementRequests()
	IncrementCache()
	IncrementExternalRequests(count int)
	GetStats() models.StatsResponse
}
