package services

import (
	"fmt"
	"weather-bot/internal/app/storage"
)

// var cityService *CityService
// var userService *UserService
// var weatherService *WeatherService
// var notificationService *NotificationService

type DualStorageError struct {
	Primary   error
	Secondary error
}

func (e *DualStorageError) Error() string {
	return fmt.Sprintf("Primary: %v, Secondary: %v", e.Primary, e.Secondary)
}

func InitCityService(primary storage.CityStorage, secondary storage.CityStorage) CityService {
	return CityService{
		Primary:   primary,
		Secondary: secondary,
	}
}

// func GetCityService() *CityService {
// 	if cityService == nil {
// 		log.Fatal().Msg("City service is not initialized")
// 	}

// 	return cityService
// }

func InitUserService(primary storage.UserStorage, secondary storage.UserStorage) UserService {
	return UserService{
		Primary:   primary,
		Secondary: secondary,
	}
}

// func GetUserService() *UserService {
// 	if userService == nil {
// 		log.Fatal().Msg("User service is not initialized")
// 	}

// 	return userService
// }

func InitWeatherService(primary storage.WeatherStorage, secondary storage.WeatherStorage) WeatherService {
	return WeatherService{
		Primary:   primary,
		Secondary: secondary,
	}

}

// func GetWeatherService() *WeatherService {
// 	if weatherService == nil {
// 		log.Fatal().Msg("Weather service is not initialized")
// 	}

// 	return weatherService
// }

func InitNotificationService(primary storage.NotificationStorage) NotificationService {
	return NotificationService{
		Primary: primary,
	}
}

// func GetNotificationService() *NotificationService {
// 	if notificationService == nil {
// 		log.Fatal().Msg("notification service is not initialized")
// 	}

// 	return notificationService
// }
