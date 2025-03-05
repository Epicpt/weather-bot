package jobs

import (
	"strconv"
	"time"
	"weather-bot/internal/app/reply"
	"weather-bot/internal/app/services"
	"weather-bot/internal/app/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func StartUserWorker(bot *tgbotapi.BotAPI) {
	time.Sleep(2 * time.Minute)
	log.Info().Msg("Воркер ProcessUserUpdate запущен...")

	for {
		ProcessUserUpdate(bot)
		log.Warn().Msg("ProcessUserUpdate завершился, перезапуск через минуту...")
		time.Sleep(1 * time.Minute)
	}
}

func ProcessUserUpdate(bot *tgbotapi.BotAPI) {
	notificationService := services.Global()
	for {
		if !notificationService.IsHealthy() {
			log.Warn().Msg("Redis недоступен, горутина ProcessUserUpdate уходит в спячку на час")
			time.Sleep(1 * time.Hour)
			continue
		}

		// Читаем задачу из `user_notifications`
		streams, err := notificationService.GetScheduleUserNotifications()

		if err != nil {
			log.Error().Err(err).Msg("Ошибка чтения уведомлений юзеров из Redis Stream")
			time.Sleep(1 * time.Minute)
			continue
		}

		// Обрабатываем задачу
		for _, stream := range streams {
			for _, message := range stream.Messages {
				userID, err := strconv.ParseInt(message.Values["user_id"].(string), 10, 64)
				if err != nil {
					log.Error().Err(err).Msg("Ошибка парсинга userID")
					continue
				}
				notifTime, _ := message.Values["executeAt"].(string)
				executeAt, err := strconv.ParseInt(notifTime, 10, 64)
				if err != nil {
					log.Error().Err(err).Msg("Ошибка парсинга executeAt")
					continue
				}
				// Если пора отправлять уведомление
				if time.Now().Unix() >= executeAt {
					log.Info().Msgf("Отправляем уведомление пользователю %d...", userID)

					user, err := services.Global().GetUser(userID)
					if err != nil {
						log.Error().Err(err).Msg("Ошибка при получении данных пользователя из Redis")
						continue
					}
					forecast, err := weather.Get(user.CityID)
					if err != nil {
						log.Error().Err(err).Msg("Ошибка при получении погоды")
						continue
					}
					today := time.Now().UTC().Format("2006-01-02")
					msg := weather.FormatDailyForecast(user.City, forecast.FullDay[today])
					reply.Send().Message(user.ChatID, msg, nil)

					notifTime := time.Unix(executeAt, 0)
					// Планируем задачу на следующий день
					ScheduleUserUpdate(userID, notifTime)
				}
			}
		}
	}

}
