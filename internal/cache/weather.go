package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"
)

var _ storage.WeatherStorage = (*Cache)(nil)

func (c *Cache) GetWeather(cityID int) (*models.ProcessedForecast, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("weather:city:%d", cityID)

	cachedData, err := c.client.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	var forecast models.ProcessedForecast
	if err := json.Unmarshal([]byte(cachedData), &forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}

func (c *Cache) SaveWeather(cityID int, forecast *models.ProcessedForecast) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("weather:city:%d", cityID)

	data, err := json.Marshal(forecast)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	err = c.client.Set(ctx, cacheKey, data, 25*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("ошибка записи в Redis: %w", err)
	}
	return nil
}
