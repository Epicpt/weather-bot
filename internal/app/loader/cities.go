package loader

import (
	"encoding/json"
	"os"
	"weather-bot/internal/app/services"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
)

func LoadCities(filePath string, service services.CityService) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var cities []models.City
	if err := json.NewDecoder(file).Decode(&cities); err != nil {
		return err
	}

	log.Info().Msgf("Загружено %d городов из файла", len(cities))

	service.LoadCities(cities)

	return nil
}
