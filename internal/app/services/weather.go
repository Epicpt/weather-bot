package services

import (
	"weather-bot/internal/app/monitoring"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
)

type WeatherService struct {
	Primary   storage.WeatherStorage
	Secondary storage.WeatherStorage
}

func (s *WeatherService) SaveWeather(id int, forecast *models.ProcessedForecast) error {
	var errP, errS error

	// Пытаемся сохранить в Primary хранилище
	errP = s.Primary.SaveWeather(id, forecast)
	if errP != nil {
		monitoring.RedisErrorsTotal.Inc()
		log.Warn().Err(errP).Int("cityID", id).Msg("Ошибка записи города в Primary хранилище")
	}

	// Пытаемся сохранить во Secondary хранилище
	errS = s.Secondary.SaveWeather(id, forecast)
	if errS != nil {
		monitoring.DBErrorsTotal.Inc()
		log.Warn().Err(errS).Int("cityID", id).Msg("Ошибка записи города в Secondary хранилище")
	}

	// Если произошли ошибки на обоих хранилищах, комбинируем ошибки и возвращаем
	if errP != nil && errS != nil {
		return &DualStorageError{Primary: errP, Secondary: errS}
	}

	// Если ошибок не было, возвращаем nil
	return nil
}

func (s *WeatherService) GetWeather(id int) (*models.ProcessedForecast, error) {
	weather, errP := s.Primary.GetWeather(id)
	if errP == nil {
		monitoring.RedisCacheHits.Inc()
		return weather, nil
	}
	monitoring.RedisCacheMisses.Inc()
	monitoring.RedisErrorsTotal.Inc()
	log.Warn().Err(errP).Msg("Ошибка чтения погоды из Primary хранилища")

	weather, errS := s.Secondary.GetWeather(id)
	if errS == nil {
		return weather, nil
	}
	monitoring.DBErrorsTotal.Inc()
	log.Warn().Err(errS).Msg("Ошибка чтения погоды из Secondary хранилища")

	return nil, &DualStorageError{Primary: errP, Secondary: errS}
}
