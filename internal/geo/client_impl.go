package geo

import "sync"

// MultiAPIClient реализация Client интерфейса
// Отправляет запросы ко всем 5 геолокационным API параллельно
type MultiAPIClient struct{}

// NewMultiAPIClient создает новый экземпляр MultiAPIClient
func NewMultiAPIClient() *MultiAPIClient {
	return &MultiAPIClient{}
}

// FetchAllAsync отправляет запросы ко всем 5 API параллельно
// и собирает результаты через канал
// Возвращает слайс результатов и количество успешных запросов
func (c *MultiAPIClient) FetchAllAsync(ip string) ([]Result, int) {
	resultsChan := make(chan Result, 5)
	var wg sync.WaitGroup

	// Запускаем 5 goroutines параллельно
	wg.Add(1)
	go func() {
		defer wg.Done()
		resultsChan <- fetchFromIPInfo(ip)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resultsChan <- fetchFromIPAPI(ip)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resultsChan <- fetchFromIPWhoIs(ip)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resultsChan <- fetchFromDBIP(ip)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resultsChan <- fetchFromIPMe(ip)
	}()

	// Закрываем канал после завершения всех запросов
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Собираем результаты
	var results []Result
	successfulCount := 0

	for result := range resultsChan {
		results = append(results, result)
		if result.Error == nil {
			successfulCount++
		}
	}

	return results, successfulCount
}
