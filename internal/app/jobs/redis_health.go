package jobs

import (
	"time"
	"weather-bot/internal/app/services"
)

func StartRedisHealthChecker() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		services.Global().HealthCheck()
	}
}
