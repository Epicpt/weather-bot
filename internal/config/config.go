package config

import (
	"os"
)

type Config struct {
	RedisURL    string
	PostgresURL string
	BotToken    string
	WeatherKey  string
}

func Load() *Config {
	// err := godotenv.Load("../internal/config/.env")
	// if err != nil {
	// 	log.Fatal("Ошибка загрузки .env файла")
	// }

	return &Config{
		RedisURL:    os.Getenv("REDIS_URL"),
		PostgresURL: os.Getenv("POSTGRES_URL"),
		BotToken:    os.Getenv("TELEGRAM_BOT_TOKEN"),
		WeatherKey:  os.Getenv("OPENWEATHER_API_KEY"),
	}
}
