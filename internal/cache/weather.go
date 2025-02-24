package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Полный прогноз на 1 день (утро, день, вечер, ночь)
type FullDayForecast struct {
	Morning WeatherSummary `json:"morning"`
	Day     WeatherSummary `json:"day"`
	Evening WeatherSummary `json:"evening"`
	Night   WeatherSummary `json:"night"`
}

// Краткая информация о погоде
type WeatherSummary struct {
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	WindSpeed   float64 `json:"wind_speed"`
	Condition   string  `json:"condition"`
	ConditionId int     `json:"conditionId"`
}

// Краткий прогноз на 5 дней (средняя температура и основное состояние погоды)
type ShortDayForecast struct {
	Date        string  `json:"date"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	ConditionId int     `json:"condition_id"`
}

// Итоговая структура, которая хранится в Redis и БД
type ProcessedForecast struct {
	FullDay   map[string]FullDayForecast `json:"full_day"`   // Прогноз на каждый день (детально)
	ShortDays []ShortDayForecast         `json:"short_days"` // Краткий прогноз на 5 дней
}

func (c *Cache) GetWeather(cityID int) (*ProcessedForecast, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("weather:city:%d", cityID)

	cachedData, err := c.c.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	var forecast ProcessedForecast
	if err := json.Unmarshal([]byte(cachedData), &forecast); err == nil {
		return nil, err
	}

	return &forecast, nil
}

func (c *Cache) SetWeather(cityID int, forecast *ProcessedForecast) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("weather:city:%d", cityID)

	data, err := json.Marshal(forecast)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	err = c.c.Set(ctx, cacheKey, data, 48*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("ошибка записи в Redis: %w", err)
	}
	return nil
}
