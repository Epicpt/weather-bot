package services

import (
	"fmt"
	"weather-bot/internal/app/monitoring"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
)

type CityService struct {
	Primary   storage.CityStorage
	Secondary storage.CityStorage
}

func (s *CityService) SaveCity(city models.City) error {
	var errP, errS error
	errP = s.Primary.SaveCity(city)
	if errP != nil {
		monitoring.RedisErrorsTotal.Inc()
		log.Warn().Err(errP).Msg("Ошибка записи города в Primary хранилище")
	}

	if errS = s.Secondary.SaveCity(city); errS != nil {
		monitoring.DBErrorsTotal.Inc()
		log.Warn().Err(errS).Msg("Ошибка записи города в Secondary хранилище")
	}

	if errP != nil && errS != nil {
		return &DualStorageError{Primary: errP, Secondary: errS}
	}
	return nil
}

func (s *CityService) GetCities(name string) ([]models.City, error) {
	cities, errP := s.Primary.GetCities(name)
	if errP == nil {
		monitoring.RedisCacheHits.Inc()
		return cities, nil
	}
	monitoring.RedisCacheMisses.Inc()
	monitoring.RedisErrorsTotal.Inc()
	log.Error().Err(errP).Msg("Ошибка получения городов из Primary хранилища")

	cities, errS := s.Secondary.GetCities(name)
	if errS == nil {
		return cities, nil
	}
	monitoring.DBErrorsTotal.Inc()
	log.Error().Err(errS).Msg("Ошибка получения городов из Secondary хранилища")
	return nil, &DualStorageError{Primary: errP, Secondary: errS}
}

func (s *CityService) LoadCities(cities []models.City) error {
	for _, city := range cities {
		if err := s.SaveCity(city); err != nil {
			return fmt.Errorf("failed to save city %s: %w", city.Name, err)
		}

	}
	return nil
}

func (s *CityService) GetCitiesNames() ([]string, error) {
	return s.getFromStorage(func(storage storage.CityStorage) ([]string, error) {
		return storage.GetCitiesNames()
	}, "имён городов")

}

func (s *CityService) GetCitiesIds() ([]string, error) {
	return s.getFromStorage(func(storage storage.CityStorage) ([]string, error) {
		return storage.GetCitiesIds()
	}, "id городов")
}

func (s *CityService) getFromStorage(get func(storage.CityStorage) ([]string, error), operation string) ([]string, error) {
	cities, errP := get(s.Primary)
	if errP == nil {
		monitoring.RedisCacheHits.Inc()
		return cities, nil
	}
	monitoring.RedisCacheMisses.Inc()
	monitoring.RedisErrorsTotal.Inc()
	log.Error().Err(errP).Msgf("Ошибка получения %s городов из Primary хранилища", operation)

	cities, errS := get(s.Secondary)
	if errS == nil {
		return cities, nil
	}
	monitoring.DBErrorsTotal.Inc()
	log.Error().Err(errS).Msgf("Ошибка получения %s городов из Secondary хранилища", operation)
	return nil, &DualStorageError{Primary: errP, Secondary: errS}
}
