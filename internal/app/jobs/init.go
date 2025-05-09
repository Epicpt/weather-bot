package jobs

import (
	"weather-bot/internal/app/monitoring"

	"github.com/rs/zerolog/log"
)

func Init() error {
	log.Info().Msg("Инициализация фоновых задач...")

	// Запуск HealthChecker
	go StartRedisHealthChecker()

	// Добавляем задачу обновления прогноза в Redis (если её нет)
	if err := ScheduleWeatherUpdate(); err != nil {
		monitoring.RedisErrorsTotal.Inc()
		log.Error().Err(err).Msg("Ошибка при установке задачи обновления погоды")
		return err
	}

	go StartWeatherWorker()
	go StartUserWorker()
	go StartCleanupTask()
	return nil
}
