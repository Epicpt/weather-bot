package jobs

import (
	"time"
	"weather-bot/internal/app/monitoring"
	"weather-bot/internal/app/services"

	"github.com/rs/zerolog/log"
)

func StartCleanupTask() {
	ticker := time.NewTicker(6 * time.Hour) // Очистка каждые 6 часов
	defer ticker.Stop()
	for {
		<-ticker.C
		err := services.Global().CleanupOldWeatherData()
		if err != nil {
			monitoring.DBErrorsTotal.Inc()
			log.Error().Err(err).Msg("Ошибка очистки данных")
		}
	}
}
