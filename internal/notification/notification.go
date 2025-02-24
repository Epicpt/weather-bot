package notification

import (
	"fmt"
	"weather-bot/internal/cache"
	"weather-bot/internal/database"

	"github.com/rs/zerolog/log"
)

func Init(redisClient *cache.Cache, db *database.Database) error {
	log.Info().Msg("Инициализация задач уведомлений...")

	// Добавляем задачу обновления прогноза в Redis (если её нет)
	if err := ScheduleWeatherUpdate(redisClient); err != nil {
		log.Error().Err(err).Msg("Ошибка при установке задачи обновления погоды")
		return err
	}

	// Запускаем воркер для обновления прогнозов в отдельной горутине
	go ProcessWeatherUpdates(redisClient, db)

	// Воркеры для пользовательских уведомлений

	return nil
}

func Get(userID int, redisClient *cache.Cache, db *database.Database) (string, error) {
	//  Пробуем получить из кеша
	timeStr, err := redisClient.GetNotificationTime(userID)
	if err == nil {
		return timeStr, nil
	}
	log.Error().Err(err).Msg("не удалось получить notification time из кеша")

	//  Если в кеше нет, берём из БД
	timeStr, err = db.GetNotificationTime(userID)
	if err != nil {
		return "", fmt.Errorf("не удалось получить notification time из БД: %w", err)
	}

	//  Обновляем кеш
	err = redisClient.SetNotificationTime(userID, timeStr)
	if err != nil {
		log.Error().Err(err).Msg("не удалось сохранить notification time в кеш")
	}

	return timeStr, nil
}

func Set(userID int, timeStr string, redisClient *cache.Cache, db *database.Database) error {
	//  Сохраняем в БД
	err := db.SetNotificationTime(userID, timeStr)
	if err != nil {
		return fmt.Errorf("не удалось сохранить notification time в БД: %w", err)
	}

	//  Обновляем кеш
	err = redisClient.SetNotificationTime(userID, timeStr)
	if err != nil {
		return fmt.Errorf("не удалось сохранить notification time в кеш")
	}

	return nil
}
