package services

import (
	"fmt"
	"weather-bot/internal/app/storage"
)

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

func InitUserService(primary storage.UserStorage, secondary storage.UserStorage) UserService {
	return UserService{
		Primary:   primary,
		Secondary: secondary,
	}
}

func InitWeatherService(primary storage.WeatherStorage, secondary storage.WeatherStorage) WeatherService {
	return WeatherService{
		Primary:   primary,
		Secondary: secondary,
	}

}

func InitNotificationService(primary storage.NotificationStorage) NotificationService {
	return NotificationService{
		Primary: primary,
	}
}
