package service

import (
	"ip-geo-checker/internal/geo"
	"ip-geo-checker/internal/models"
	"ip-geo-checker/internal/repository"
	"log"
	"net"
)

// IPCheckerService сервис для проверки IP адресов
// Содержит всю бизнес-логику приложения
type IPCheckerService struct {
	cacheRepo repository.CacheRepository
	statsRepo repository.StatsRepository
	geoClient geo.Client
}

// NewIPCheckerService создает новый экземпляр IPCheckerService
func NewIPCheckerService(
	cacheRepo repository.CacheRepository,
	statsRepo repository.StatsRepository,
	geoClient geo.Client,
) *IPCheckerService {
	return &IPCheckerService{
		cacheRepo: cacheRepo,
		statsRepo: statsRepo,
		geoClient: geoClient,
	}
}

// ValidateIP проверяет валидность IP адреса
func (s *IPCheckerService) ValidateIP(ip string) error {
	if ip == "" {
		return ErrEmptyIP
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ErrInvalidIPFormat
	}

	if parsedIP.String() == "" {
		return ErrInvalidIPFormat
	}

	return nil
}

// CheckIP проверяет IP адрес и возвращает информацию о городе
func (s *IPCheckerService) CheckIP(ip string) (models.CheckIPResponse, error) {
	// Проверяем кэш
	if cached, exists := s.cacheRepo.Get(ip); exists {
		log.Printf("Cache hit for IP %s", ip)
		s.statsRepo.IncrementCache()
		cached.FromCache = true
		return cached, nil
	}

	// Данных нет в кэше - делаем запросы к внешним API
	log.Printf("Cache miss for IP %s, fetching from external APIs", ip)
	results, successfulCount := s.geoClient.FetchAllAsync(ip)

	// Обновляем статистику успешных внешних запросов
	s.statsRepo.IncrementExternalRequests(successfulCount)

	// Подсчитываем проценты городов
	cityPercentages := calculateCityPercentages(results)

	// Если не удалось получить данные ни от одного API
	if cityPercentages == "" {
		log.Printf("Warning: Failed to get city information for IP %s from any API", ip)
		return models.CheckIPResponse{}, ErrNoDataFromAPIs
	}

	log.Printf("Successfully retrieved city information for IP %s: %s", ip, cityPercentages)

	// Формируем ответ
	response := models.CheckIPResponse{
		City:      cityPercentages,
		FromCache: false,
	}

	// Сохраняем в кэш
	s.cacheRepo.Set(ip, response)

	return response, nil
}

// IncrementRequests увеличивает счетчик запросов
func (s *IPCheckerService) IncrementRequests() {
	s.statsRepo.IncrementRequests()
}

// GetStats возвращает статистику приложения
func (s *IPCheckerService) GetStats() models.StatsResponse {
	stats := s.statsRepo.GetStats()
	stats.CacheCount = s.cacheRepo.Count()
	return stats
}
