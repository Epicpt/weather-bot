package services

import (
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
)

type WeatherService struct {
	Primary   storage.WeatherStorage
	Secondary storage.WeatherStorage
}

func (s *WeatherService) SaveWeather(id int, forecast *models.ProcessedForecast) error {
	errP := s.Primary.SaveWeather(id, forecast)
	if errP != nil {
		log.Warn().Err(errP).Msg("Ошибка записи города в Primary хранилище")
	}
	if errS := s.Secondary.SaveWeather(id, forecast); errS != nil {
		log.Warn().Err(errS).Msg("Ошибка записи города в Secondary хранилище")
		return &DualStorageError{Primary: errP, Secondary: errS}
	}
	return nil
}

func (s *WeatherService) GetWeather(id int) (*models.ProcessedForecast, error) {
	weather, errP := s.Primary.GetWeather(id)
	if errP == nil {
		return weather, nil
	}
	log.Warn().Err(errP).Msg("Ошибка чтения погоды из Primary хранилища")

	weather, errS := s.Secondary.GetWeather(id)
	if errS == nil {
		return weather, nil
	}
	log.Warn().Err(errS).Msg("Ошибка чтения погоды из Secondary хранилища")

	return nil, &DualStorageError{Primary: errP, Secondary: errS}
}
