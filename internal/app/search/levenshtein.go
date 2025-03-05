package search

import (
	"fmt"
	"sort"
	"weather-bot/internal/app/services"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// FindTop3ClosestCities находит 3 похожих города
func findTop3ClosestCities(input string) ([]models.City, error) {

	type cityDistance struct {
		city     string
		distance int
	}

	var distances []cityDistance

	cityNames, err := services.Global().GetCitiesNames()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения имен городов: %w", err)
	}

	// Считаем расстояния
	for _, city := range cityNames {
		distance := levenshtein.DistanceForStrings([]rune(input), []rune(city), levenshtein.DefaultOptions)
		if distance <= 7 {
			distances = append(distances, cityDistance{city: city, distance: distance})
		}

	}

	if len(distances) == 0 {
		return nil, fmt.Errorf("похожие города не найдены")
	}

	// Сортируем по возрастанию расстояния
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	var closestCities []models.City

	for i := 0; i < 3 && i < len(distances); i++ {
		citiesClose, err := services.Global().GetCities(distances[i].city)
		if err != nil {
			log.Error().Err(err).Msgf("Ошибка при получении города %s", distances[i].city)
			continue
		}

		closestCities = append(closestCities, citiesClose...)

	}

	if len(closestCities) == 0 {
		return nil, fmt.Errorf("похожие города не найдены")
	}

	return closestCities, nil
}
