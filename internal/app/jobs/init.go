package jobs

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func Init(bot *tgbotapi.BotAPI) error {
	log.Info().Msg("Инициализация фоновых задач...")

	// Запуск HealthChecker
	go StartRedisHealthChecker()

	// Добавляем задачу обновления прогноза в Redis (если её нет)
	if err := ScheduleWeatherUpdate(); err != nil {
		log.Error().Err(err).Msg("Ошибка при установке задачи обновления погоды")
		return err
	}

	go StartWeatherWorker()
	go StartUserWorker(bot)
	return nil
}
