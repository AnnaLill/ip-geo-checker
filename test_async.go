package main

import (
	"fmt"
)

// testAsyncRequests тестирует асинхронные запросы
// Этот файл можно запустить отдельно для проверки работы функции
func testAsyncRequests() {
	fmt.Println("Testing async requests to all 5 APIs...")
	fmt.Println("IP: 8.8.8.8")

	results, successfulCount := fetchAllAPIsAsync("8.8.8.8")

	fmt.Printf("\nTotal results: %d\n", len(results))
	fmt.Printf("Successful requests: %d\n", successfulCount)
	fmt.Println("\nDetailed results:")

	for i, result := range results {
		// Используем имя API из результата, если оно есть
		apiName := result.APIName
		if apiName == "" {
			apiName = fmt.Sprintf("API #%d", i+1)
		}
		fmt.Printf("\n%d. %s:\n", i+1, apiName)

		if result.Error != nil {
			fmt.Printf("   Error: %v\n", result.Error)
		} else {
			fmt.Printf("   City: %s\n", result.City)
		}
	}
}

// Для запуска теста используйте команду:
// PowerShell: go run test_async.go parsers.go
// Или: cd в директорию проекта, затем: go run test_async.go parsers.go

/*func main() {
	testAsyncRequests()
}
*/
