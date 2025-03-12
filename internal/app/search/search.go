package search

import (
	"fmt"
	"weather-bot/internal/app/services"
	"weather-bot/internal/models"
	"weather-bot/pkg/utils"

	"github.com/rs/zerolog/log"
)

// SearchCity ищет город в хранилищах и похожие на ввод
func SearchCity(cityName string) ([]models.City, error) {
	cityName = utils.NormalizeCityName(cityName)

	cities, err := services.Global().GetCities(cityName)
	if err != nil {
		log.Debug().Err(err).Msg("Ошибка получения городов из хранилищ")
	}

	if cities == nil || len(cities) == 0 {
		closestMatch, err := findTop3ClosestCities(cityName)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения похожих городов: %w", err)
		}
		if closestMatch == nil {
			return nil, fmt.Errorf("Похожие города с таким названием не найдены")
		}
		return closestMatch, nil
	}

	return cities, nil

}
