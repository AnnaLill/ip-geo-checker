package main

import (
	"fmt"
	"sort"
	"strings"
)

// calculateCityPercentages подсчитывает проценты для каждого города
// из результатов всех API и формирует строку в формате "City1 (X%), City2 (Y%)"
func calculateCityPercentages(results []GeoAPIResult) string {
	// Создаем map для подсчета количества упоминаний каждого города
	cityCounts := make(map[string]int)
	totalSuccessful := 0

	// Проходим по всем результатам и считаем города
	for _, result := range results {
		// Учитываем только успешные запросы (без ошибок)
		if result.Error == nil && result.City != "" {
			cityCounts[result.City]++
			totalSuccessful++
		}
	}

	// Если нет успешных результатов, возвращаем пустую строку
	if totalSuccessful == 0 {
		return ""
	}

	// Создаем слайс структур для сортировки
	type cityPercent struct {
		city    string
		count   int
		percent int
	}

	var cities []cityPercent
	for city, count := range cityCounts {
		percent := (count * 100) / totalSuccessful
		cities = append(cities, cityPercent{
			city:    city,
			count:   count,
			percent: percent,
		})
	}

	// Сортируем по количеству упоминаний (от большего к меньшему)
	sort.Slice(cities, func(i, j int) bool {
		if cities[i].count != cities[j].count {
			return cities[i].count > cities[j].count
		}
		return cities[i].city < cities[j].city
	})

	// Формируем строку результата
	var parts []string
	for _, cp := range cities {
		parts = append(parts, fmt.Sprintf("%s (%d%%)", cp.city, cp.percent))
	}

	return strings.Join(parts, ", ")
}
