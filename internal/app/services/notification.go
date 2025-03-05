package services

import (
	"weather-bot/internal/app/storage"

	"github.com/redis/go-redis/v9"
)

type NotificationService struct {
	Primary storage.NotificationStorage
}

func (s *NotificationService) GetUserNotificationTime(userID int64) (string, error) {
	return s.Primary.GetUserNotificationTime(userID)
}

func (s *NotificationService) RemoveUserNotification(userID int64) error {
	return s.Primary.RemoveUserNotification(userID)
}

func (s *NotificationService) ScheduleUserNotification(userID int64, executeAt int64) error {
	return s.Primary.ScheduleUserNotification(userID, executeAt)
}

func (s *NotificationService) ScheduleWeatherUpdate(executeAt int64) error {
	return s.Primary.ScheduleWeatherUpdate(executeAt)
}

func (s *NotificationService) RemoveWeatherUpdate() error {
	return s.Primary.RemoveWeatherUpdate()
}

func (s *NotificationService) GetScheduleWeatherUpdate() ([]redis.XStream, error) {
	return s.Primary.GetScheduleWeatherUpdate()
}
func (s *NotificationService) GetScheduleUserNotifications() ([]redis.XStream, error) {
	return s.Primary.GetScheduleUserNotifications()
}
