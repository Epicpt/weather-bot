package main

import (
	"weather-bot/internal/app"
	"weather-bot/internal/config"
)

func main() {
	cfg := config.Load()

	var aplication app.Application = app.New(cfg)

	aplication.Bootstrap()
	defer aplication.Shutdown()

	aplication.Run()
}
