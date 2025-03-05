package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
)

var _ storage.CityStorage = (*Cache)(nil)

func (c *Cache) SaveCity(city models.City) error {
	redisKey := fmt.Sprintf("city:%s", city.Name)

	// Преобразуем структуру в JSON
	cityData, err := json.Marshal(city)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при сериализации города")
		return fmt.Errorf("ошибка при сериализации города: %w", err)
	}

	// Проверяем, существует ли уже город с таким названием
	existingCities, err := c.client.LRange(context.Background(), redisKey, 0, -1).Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при проверке существующих городов")
		return fmt.Errorf("ошибка при проверке существующих городов: %w", err)
	}

	for i, existingCityData := range existingCities {
		var existingCity models.City
		if err := json.Unmarshal([]byte(existingCityData), &existingCity); err != nil {
			continue
		}
		if existingCity.ID == city.ID {
			// Город с таким ID уже существует, обновляем его
			err = c.client.LSet(context.Background(), redisKey, int64(i), string(cityData)).Err()
			if err != nil {
				log.Error().Err(err).Msg("Ошибка при обновлении города в Redis")
				return fmt.Errorf("ошибка при обновлении города в Redis: %w", err)
			}
			return nil
		}
	}

	// Если город не найден, добавляем новый
	err = c.client.RPush(context.Background(), redisKey, cityData).Err()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка записи в Redis")
		return fmt.Errorf("ошибка записи в Redis: %w", err)
	}

	return nil
}

func (c *Cache) GetCities(city string) ([]models.City, error) {
	seen := make(map[string]bool)
	var result []models.City
	redisKey := fmt.Sprintf("city:%s", city)

	citiesData, err := c.client.LRange(context.Background(), redisKey, 0, -1).Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка получения городов из Redis")
		return nil, fmt.Errorf("ошибка получения городов из Redis: %w", err)
	}
	for _, data := range citiesData {
		var city models.City
		err := json.Unmarshal([]byte(data), &city)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка десериализации города")
			continue
		}
		if seen[city.Region] {
			continue
		}

		seen[city.Region] = true
		result = append(result, city)
	}
	return result, nil
}

func (c *Cache) GetCitiesNames() ([]string, error) {

	// Получаем все ключи, соответствующие шаблону "city:*"
	keys, err := c.client.Keys(context.Background(), "city:*").Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ключей из Redis: %w", err)
	}

	// Удаляем префикс "city:" из каждого ключа
	cities := make([]string, 0, len(keys))
	for _, key := range keys {
		cityName := strings.TrimPrefix(key, "city:")
		cities = append(cities, cityName)
	}

	return cities, nil
}

func (c *Cache) GetCitiesIds() ([]string, error) {
	ctx := context.Background()

	// Получаем все ключи, соответствующие шаблону "user:*"
	userKeys, err := c.client.Keys(ctx, "user:*").Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка получения user-ключей из Redis")
		return nil, err
	}

	var cityIDs []string
	seenCities := make(map[string]bool)

	// Извлекаем `city_id` у каждого пользователя
	for _, key := range userKeys {
		cityID, err := c.client.HGet(ctx, key, "city_id").Result()
		if err == nil && !seenCities[cityID] {
			cityIDs = append(cityIDs, cityID)
			seenCities[cityID] = true
		}
	}

	if len(cityIDs) == 0 {
		return nil, fmt.Errorf("нет данных в Redis, нужно запросить из БД")
	}

	return cityIDs, nil
}
