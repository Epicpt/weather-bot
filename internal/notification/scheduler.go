package notification

import (
	"fmt"
	"strconv"
	"time"
	"weather-bot/internal/cache"
	"weather-bot/internal/database"
	"weather-bot/internal/weather"

	"github.com/rs/zerolog/log"
)

func ScheduleWeatherUpdate(redisClient *cache.Cache) error {
	// Вычисляем `executeAt` (00:01 следующего дня)
	now := time.Now()
	executeAt := time.Date(now.Year(), now.Month(), now.Day(), 0, 1, 0, 0, now.Location()).Add(24 * time.Hour).UnixMilli()

	err := redisClient.SetUpdateWeather(executeAt)
	if err != nil {
		return fmt.Errorf("не удалось сохранить executeAt в Redis: %w", err)
	}

	log.Info().Msgf("Задача на обновление погоды запланирована на %s", time.UnixMilli(executeAt).Format("15:04"))
	return nil
}

func ProcessWeatherUpdates(redisClient *cache.Cache, db *database.Database) {

	for {
		// Читаем задачу из `weather_updates`
		streams, err := redisClient.GetUpdateWeather()

		if err != nil {
			log.Error().Err(err).Msg("Ошибка чтения задачи обновления погоды из Redis Stream")
			continue
		}

		// Обрабатываем задачу
		for _, stream := range streams {
			for _, message := range stream.Messages {
				executeAt, _ := strconv.ParseInt(message.Values["executeAt"].(string), 10, 64)

				// Если время выполнения уже пришло
				if time.Now().UnixMilli() >= executeAt {
					log.Info().Msg("Запуск обновления погоды...")

					cityIDs, err := redisClient.GetAllCitiesIds()
					if err != nil {
						log.Error().Err(err).Msg("Ошибка получения городов из redis")
						cityIDs, err = db.GetCitiesIds()
						if err != nil {
							log.Error().Err(err).Msg("Ошибка получения городов из БД")
							continue
						}
					}
					err = weather.Update(cityIDs, redisClient, db)
					if err != nil {
						log.Error().Err(err).Msg("Ошибка при обновлении погоды")
					} else {
						log.Info().Msg("Погода успешно обновлена")
					}

					// Удаляем задачу
					err = redisClient.DeleteUpdateWeather()
					if err != nil {
						log.Error().Err(err).Msg("Ошибка удаления задачи обновления погоды из Redis")
					}

					// Планируем задачу на следующий день
					ScheduleWeatherUpdate(redisClient)
				}
			}
		}
	}

}
