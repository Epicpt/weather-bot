package jobs

import (
	"fmt"
	"time"
	"weather-bot/internal/app/services"

	"github.com/rs/zerolog/log"
)

func ScheduleWeatherUpdate() error {
	notificationService := services.Global()
	// Удаляем задачу
	err := notificationService.RemoveWeatherUpdate()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка удаления задачи обновления погоды из Redis")
	}

	// Вычисляем `executeAt` (00:01 следующего дня)
	now := time.Now()
	executeAt := time.Date(now.Year(), now.Month(), now.Day(), 0, 1, 0, 0, now.Location()).Add(24 * time.Hour).Unix()

	err = notificationService.ScheduleWeatherUpdate(executeAt)
	if err != nil {
		return fmt.Errorf("не удалось сохранить executeAt в Redis: %w", err)
	}

	log.Info().Msgf("Задача на обновление погоды запланирована на %s", time.Unix(executeAt, 0).Format("15:04"))
	return nil
}

func ScheduleUserUpdate(userID int64, notificationTime time.Time) error {
	notificationService := services.Global()
	// Удаляем задачу
	err := notificationService.RemoveUserNotification(userID)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка удаления задачи обновления погоды из Redis")
	}

	// Вычисляем время следующего уведомления
	now := time.Now()
	executeAt := time.Date(now.Year(), now.Month(), now.Day(), notificationTime.Hour(), notificationTime.Minute(), 0, 0, now.Location()).Add(24 * time.Hour).Unix()

	// Сохраняем новую задачу
	err = notificationService.ScheduleUserNotification(userID, executeAt)
	if err != nil {
		return fmt.Errorf("не удалось сохранить executeAt в Redis: %w", err)
	}

	log.Info().Msgf("Задача на обновление погоды для юзера %d запланирована на %s", userID, time.Unix(executeAt, 0).Format("15:04"))
	return nil
}
