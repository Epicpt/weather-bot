package storage

import (
	"weather-bot/internal/models"

	"github.com/redis/go-redis/v9"
)

type Cashe interface {
	DB
	NotificationStorage
	HealthChecker
}
type DB interface {
	CityStorage
	UserStorage
	WeatherStorage
}

type CityStorage interface {
	SaveCity(models.City) error
	GetCities(string) ([]models.City, error)
	GetCitiesNames() ([]string, error)
	GetCitiesIds() ([]string, error)
}

type UserStorage interface {
	SaveUser(*models.User) error
	GetUser(int64) (*models.User, error)
}

type WeatherStorage interface {
	SaveWeather(int, *models.ProcessedForecast) error
	GetWeather(int) (*models.ProcessedForecast, error)
}

type NotificationStorage interface {
	GetUserNotificationTime(int64) (string, error)
	RemoveUserNotification(int64) error
	ScheduleUserNotification(int64, int64) error
	ScheduleWeatherUpdate(int64) error
	RemoveWeatherUpdate() error
	GetScheduleWeatherUpdate() ([]redis.XStream, error)
	GetScheduleUserNotifications() ([]redis.XStream, error)
}

type HealthChecker interface {
	HealthCheck()
	IsHealthy() bool
}
