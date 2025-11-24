package service

import (
	"fmt"
	"ip-geo-checker/internal/geo"
	"sort"
	"strings"
)

// calculateCityPercentages подсчитывает проценты для каждого города
// из результатов всех API и формирует строку в формате "City1 (X%), City2 (Y%)"
func calculateCityPercentages(results []geo.Result) string {
	cityCounts := make(map[string]int)
	totalSuccessful := 0

	for _, result := range results {
		if result.Error == nil && result.City != "" {
			cityCounts[result.City]++
			totalSuccessful++
		}
	}

	if totalSuccessful == 0 {
		return ""
	}

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

	sort.Slice(cities, func(i, j int) bool {
		if cities[i].count != cities[j].count {
			return cities[i].count > cities[j].count
		}
		return cities[i].city < cities[j].city
	})

	var parts []string
	for _, cp := range cities {
		parts = append(parts, fmt.Sprintf("%s (%d%%)", cp.city, cp.percent))
	}

	return strings.Join(parts, ", ")
}
