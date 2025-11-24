package repository

import "ip-geo-checker/internal/models"

// CacheRepository интерфейс для работы с кэшем
// Использование интерфейса позволяет легко заменить реализацию (например, на Redis)
type CacheRepository interface {
	Get(ip string) (models.CheckIPResponse, bool)
	Set(ip string, response models.CheckIPResponse)
	Count() int
}
