package database

import (
	"context"
	"encoding/json"
	"fmt"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"
)

var _ storage.WeatherStorage = (*Database)(nil)

func (d *Database) GetWeather(cityID int) (*models.ProcessedForecast, error) {
	var forecastJSON string
	err := d.pool.QueryRow(context.Background(), "SELECT forecast FROM weather WHERE city_id = $1", cityID).Scan(&forecastJSON)
	if err != nil {
		return nil, err
	}

	var forecast models.ProcessedForecast
	if err := json.Unmarshal([]byte(forecastJSON), &forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}

func (d *Database) SaveWeather(cityID int, forecast *models.ProcessedForecast) error {
	data, err := json.Marshal(forecast)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %w", err)
	}
	_, err = d.pool.Exec(context.Background(), `
    INSERT INTO weather (city_id, forecast, updated_at) 
    VALUES ($1, $2, NOW()) 
    ON CONFLICT (city_id) DO UPDATE 
    SET forecast = $2, updated_at = NOW()`, cityID, data)
	if err != nil {
		return fmt.Errorf("ошибка сохранения в БД: %w", err)
	}
	return nil
}

func (d *Database) CleanupOldWeatherData() error {
	_, err := d.pool.Exec(context.Background(), "DELETE FROM weather WHERE updated_at < NOW() - INTERVAL '2 days'")
	return err
}
