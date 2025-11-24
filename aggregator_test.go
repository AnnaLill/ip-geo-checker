package main

import (
	"fmt"
	"testing"
)

// TestCalculateCityPercentages - тест для проверки подсчета процентов
func TestCalculateCityPercentages(t *testing.T) {
	// Тестовые данные: 3 раза Moscow, 2 раза SPB
	results := []GeoAPIResult{
		{APIName: "api1", City: "Moscow", Error: nil},
		{APIName: "api2", City: "Moscow", Error: nil},
		{APIName: "api3", City: "Moscow", Error: nil},
		{APIName: "api4", City: "SPB", Error: nil},
		{APIName: "api5", City: "SPB", Error: nil},
	}

	result := calculateCityPercentages(results)
	expected := "Moscow (60%), SPB (40%)"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestCalculateCityPercentagesWithErrors - тест с ошибками
func TestCalculateCityPercentagesWithErrors(t *testing.T) {
	// Тестовые данные: 2 успешных, 3 ошибки
	results := []GeoAPIResult{
		{APIName: "api1", City: "Moscow", Error: nil},
		{APIName: "api2", City: "Moscow", Error: nil},
		{APIName: "api3", City: "", Error: nil},                 // Пустой город
		{APIName: "api4", City: "", Error: fmt.Errorf("error")}, // Ошибка
		{APIName: "api5", City: "", Error: fmt.Errorf("error")}, // Ошибка
	}

	result := calculateCityPercentages(results)
	expected := "Moscow (100%)"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestCalculateCityPercentagesEmpty - тест с пустыми результатами
func TestCalculateCityPercentagesEmpty(t *testing.T) {
	results := []GeoAPIResult{
		{APIName: "api1", City: "", Error: fmt.Errorf("error")},
		{APIName: "api2", City: "", Error: fmt.Errorf("error")},
	}

	result := calculateCityPercentages(results)
	expected := ""

	if result != expected {
		t.Errorf("Expected empty string, got %s", result)
	}
}
