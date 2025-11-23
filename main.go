package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// Request структура для входящего запроса POST /check_ip
type CheckIPRequest struct {
	IP string `json:"ip"`
}

// Response структура для ответа POST /check_ip
type CheckIPResponse struct {
	City      string `json:"city"`
	FromCache bool   `json:"from_cache"`
}

// StatsResponse структура для ответа GET /get_stats
type StatsResponse struct {
	RequestsCount                   int `json:"requests_count"`
	CacheCount                      int `json:"cache_count"`
	ExternalSuccessfulRequestsCount int `json:"external_successful_requests_count"`
}

// Cache структура для хранения кэшированных данных
type Cache struct {
	mu    sync.RWMutex
	items map[string]CheckIPResponse
}

// Stats структура для хранения статистики
type Stats struct {
	mu                              sync.RWMutex
	requestsCount                   int
	cacheCount                      int
	externalSuccessfulRequestsCount int
}

// NewCache создает новый экземпляр кэша
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CheckIPResponse),
	}
}

// NewStats создает новый экземпляр статистики
func NewStats() *Stats {
	return &Stats{}
}

// Get получает значение из кэша
func (c *Cache) Get(ip string) (CheckIPResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.items[ip]
	return value, exists
}

// Set сохраняет значение в кэш
func (c *Cache) Set(ip string, response CheckIPResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[ip] = response
}

// Count возвращает количество элементов в кэше
func (c *Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// IncrementRequests увеличивает счетчик запросов
func (s *Stats) IncrementRequests() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestsCount++
}

// IncrementCache увеличивает счетчик кэша
func (s *Stats) IncrementCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cacheCount++
}

// IncrementExternalRequests увеличивает счетчик успешных внешних запросов
func (s *Stats) IncrementExternalRequests(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.externalSuccessfulRequestsCount += count
}

// GetStats возвращает текущую статистику
func (s *Stats) GetStats() StatsResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return StatsResponse{
		RequestsCount:                   s.requestsCount,
		CacheCount:                      s.cacheCount,
		ExternalSuccessfulRequestsCount: s.externalSuccessfulRequestsCount,
	}
}

// App структура приложения, хранит зависимости
type App struct {
	cache *Cache
	stats *Stats
}

// NewApp создает новый экземпляр приложения
func NewApp() *App {
	return &App{
		cache: NewCache(),
		stats: NewStats(),
	}
}

// handleCheckIP обрабатывает POST /check_ip
func (app *App) handleCheckIP(w http.ResponseWriter, r *http.Request) {
	// Увеличиваем счетчик запросов
	app.stats.IncrementRequests()

	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON из тела запроса
	var req CheckIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// TODO: Здесь будет логика проверки IP и запросов к внешним API
	// Пока возвращаем заглушку
	response := CheckIPResponse{
		City:      "Not implemented yet",
		FromCache: false,
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetStats обрабатывает GET /get_stats
func (app *App) handleGetStats(w http.ResponseWriter, r *http.Request) {
	// Увеличиваем счетчик запросов
	app.stats.IncrementRequests()

	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем статистику
	stats := app.stats.GetStats()
	stats.CacheCount = app.cache.Count()

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	// Создаем экземпляр приложения
	app := NewApp()

	// Настраиваем роутинг
	http.HandleFunc("/check_ip", app.handleCheckIP)
	http.HandleFunc("/get_stats", app.handleGetStats)

	// Запускаем сервер на порту 8080
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
