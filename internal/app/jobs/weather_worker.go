package jobs

import (
	"strconv"
	"time"
	"weather-bot/internal/app/monitoring"
	"weather-bot/internal/app/services"
	"weather-bot/internal/app/weather"

	"github.com/rs/zerolog/log"
)

func StartWeatherWorker() {
	time.Sleep(2 * time.Minute)
	log.Info().Msg("Воркер ProcessWeatherUpdates запущен...")

	for {
		ProcessWeatherUpdates()
		log.Warn().Msg("ProcessWeatherUpdates завершился, перезапуск через минуту...")
		time.Sleep(1 * time.Minute)
	}
}

func ProcessWeatherUpdates() {
	notificationService := services.Global()
	for {
		if !notificationService.IsHealthy() {
			log.Warn().Msg("Redis недоступен, горутина ProcessWeatherUpdates уходит в спячку на час")
			time.Sleep(1 * time.Hour)
			continue
		}
		// Читаем задачу из `weather_updates`
		streams, err := notificationService.GetScheduleWeatherUpdate()

		if err != nil {
			monitoring.RedisErrorsTotal.Inc()
			log.Error().Err(err).Msg("Ошибка чтения задачи обновления погоды из Redis Stream")
			time.Sleep(1 * time.Minute)
			continue
		}

		// Обрабатываем задачу
		for _, stream := range streams {
			for _, message := range stream.Messages {
				executeAt, err := strconv.ParseInt(message.Values["executeAt"].(string), 10, 64)
				if err != nil {
					monitoring.WeatherUpdateFailed.Inc()
					log.Error().Err(err).Int64("executeAt", executeAt).Msg("Ошибка парсинга executeAt")
					continue
				}

				// Если время выполнения уже пришло
				if time.Now().Unix() >= executeAt {
					log.Info().Msg("Запуск обновления погоды...")

					cityIDs, err := services.Global().GetCitiesIds()
					if err != nil {
						monitoring.WeatherUpdateFailed.Inc()
						log.Error().Err(err).Msg("Ошибка получения городов из хранилищ")
						continue
					}
					err = weather.Update(cityIDs)
					if err != nil {
						monitoring.WeatherUpdateFailed.Inc()
						log.Error().Err(err).Msg("Ошибка при обновлении погоды")
					} else {
						log.Info().Msg("Погода успешно обновлена")
					}
					monitoring.WeatherUpdateTotal.Inc()

					// Планируем задачу на следующий день
					ScheduleWeatherUpdate()
				}
			}
		}
	}

}
