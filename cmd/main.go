package main

import (
	"weather-bot/internal/app"
	"weather-bot/internal/config"
	"weather-bot/pkg/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	// Инициализация логгера
	logger.New()
	log.Info().Msg("Logger initialized")

	cfg := config.Load()
	log.Info().Msg("Config initialized")

	var aplication app.Application = app.New(cfg)

	aplication.Bootstrap()
	defer aplication.Shutdown()

	aplication.Run()
}
