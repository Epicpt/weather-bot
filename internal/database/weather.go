package database

import (
	"context"
	"encoding/json"
	"fmt"

	"weather-bot/internal/cache"
)

func (d *Database) GetWeather(cityID int) (*cache.ProcessedForecast, error) {
	var forecastJSON string
	err := d.db.QueryRow(context.Background(), "SELECT forecast FROM weather WHERE city_id = $1", cityID).Scan(&forecastJSON)
	if err != nil {
		return nil, err
	}

	var forecast cache.ProcessedForecast
	if err := json.Unmarshal([]byte(forecastJSON), &forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}

func (d *Database) SetWeather(cityID int, forecast *cache.ProcessedForecast) error {
	data, err := json.Marshal(forecast)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %w", err)
	}
	_, err = d.db.Exec(context.Background(), "INSERT INTO weather (city_id, forecast) VALUES ($1, $2) ON CONFLICT (city_id) DO UPDATE SET forecast = $2", cityID, data)
	if err != nil {
		return fmt.Errorf("ошибка сохранения в БД: %w", err)
	}
	return nil
}
