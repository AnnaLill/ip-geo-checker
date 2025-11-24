package handler

import (
	"encoding/json"
	"ip-geo-checker/internal/models"
	"ip-geo-checker/internal/service"
	"log"
	"net/http"
)

// Handler структура для HTTP обработчиков
type Handler struct {
	service *service.IPCheckerService
}

// NewHandler создает новый экземпляр Handler
func NewHandler(service *service.IPCheckerService) *Handler {
	return &Handler{
		service: service,
	}
}

// HandleCheckIP обрабатывает POST /check_ip
func (h *Handler) HandleCheckIP(w http.ResponseWriter, r *http.Request) {
	// Увеличиваем счетчик запросов
	h.service.IncrementRequests()

	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON из тела запроса
	var req models.CheckIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидируем IP
	if err := h.service.ValidateIP(req.IP); err != nil {
		log.Printf("Invalid IP address: %s, error: %v", req.IP, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем данные о городе
	response, err := h.service.CheckIP(req.IP)
	if err != nil {
		if err == service.ErrNoDataFromAPIs {
			http.Error(w, "Unable to get city information from any API. Please try again later.", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleGetStats обрабатывает GET /get_stats
func (h *Handler) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	// Увеличиваем счетчик запросов
	h.service.IncrementRequests()

	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем статистику
	stats := h.service.GetStats()

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding stats response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
