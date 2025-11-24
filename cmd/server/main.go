package main

import (
	"ip-geo-checker/internal/geo"
	"ip-geo-checker/internal/handler"
	"ip-geo-checker/internal/repository"
	"ip-geo-checker/internal/service"
	"log"
	"net/http"
)

func main() {
	// Инициализация зависимостей (Dependency Injection)
	// Repository слой
	cacheRepo := repository.NewInMemoryCache()
	statsRepo := repository.NewInMemoryStats()

	// Geo клиент
	geoClient := geo.NewMultiAPIClient()

	// Service слой
	ipCheckerService := service.NewIPCheckerService(cacheRepo, statsRepo, geoClient)

	// Handler слой
	h := handler.NewHandler(ipCheckerService)

	// Настраиваем роутинг
	http.HandleFunc("/check_ip", h.HandleCheckIP)
	http.HandleFunc("/get_stats", h.HandleGetStats)

	// Запускаем сервер на порту 8080
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
