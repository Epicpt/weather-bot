package services

import (
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/redis/go-redis/v9"
)

var globalStorage *ServiceContainer

type ServiceContainer struct {
	CityService         CityService
	UserService         UserService
	WeatherService      WeatherService
	NotificationService NotificationService
	Cache               storage.Cache
	DB                  storage.Database
}

func Init(primary storage.Cache, secondary storage.Database) {
	globalStorage = &ServiceContainer{
		CityService:         InitCityService(primary, secondary),
		UserService:         InitUserService(primary, secondary),
		WeatherService:      InitWeatherService(primary, secondary),
		NotificationService: InitNotificationService(primary),
		Cache:               primary,
		DB:                  secondary,
	}
}

func Global() *ServiceContainer {
	return globalStorage
}

func (s *ServiceContainer) HealthCheck() {
	s.Cache.HealthCheck()
}

func (s *ServiceContainer) IsHealthy() bool {
	return s.Cache.IsHealthy()
}

func (s *ServiceContainer) CleanupOldWeatherData() error {
	return s.DB.CleanupOldWeatherData()
}

func (s *ServiceContainer) SaveCity(city models.City) error {
	return s.CityService.SaveCity(city)
}

func (s *ServiceContainer) GetCities(city string) ([]models.City, error) {
	return s.CityService.GetCities(city)
}

func (s *ServiceContainer) GetCitiesNames() ([]string, error) {
	return s.CityService.GetCitiesNames()

}

func (s *ServiceContainer) GetCitiesIds() ([]string, error) {
	return s.CityService.GetCitiesIds()
}

func (s *ServiceContainer) LoadCities(cities []models.City) {
	s.CityService.LoadCities(cities)
}

func (s *ServiceContainer) GetUserNotificationTime(userID int64) (string, error) {
	return s.NotificationService.GetUserNotificationTime(userID)
}
func (s *ServiceContainer) RemoveUserNotification(userID int64) error {
	return s.NotificationService.RemoveUserNotification(userID)
}

func (s *ServiceContainer) ScheduleUserNotification(userID int64, executeAt int64) error {
	return s.NotificationService.ScheduleUserNotification(userID, executeAt)
}

func (s *ServiceContainer) ScheduleWeatherUpdate(executeAt int64) error {
	return s.NotificationService.ScheduleWeatherUpdate(executeAt)
}

func (s *ServiceContainer) RemoveWeatherUpdate() error {
	return s.NotificationService.RemoveWeatherUpdate()
}

func (s *ServiceContainer) GetScheduleWeatherUpdate() ([]redis.XStream, error) {
	return s.NotificationService.GetScheduleWeatherUpdate()
}
func (s *ServiceContainer) GetScheduleUserNotifications() ([]redis.XStream, error) {
	return s.NotificationService.GetScheduleUserNotifications()
}

func (s *ServiceContainer) SaveUser(user *models.User) error {
	return s.UserService.SaveUser(user)
}

func (s *ServiceContainer) GetUser(id int64) (*models.User, error) {
	return s.UserService.GetUser(id)
}

func (s *ServiceContainer) SaveWeather(id int, forecast *models.ProcessedForecast) error {
	return s.WeatherService.SaveWeather(id, forecast)
}

func (s *ServiceContainer) GetWeather(id int) (*models.ProcessedForecast, error) {
	return s.WeatherService.GetWeather(id)
}
