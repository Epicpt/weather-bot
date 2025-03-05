package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr   string
	PostgresURL string
	BotToken    string
	WeatherKey  string
}

func Load() *Config {
	err := godotenv.Load("../internal/config/.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	return &Config{
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		PostgresURL: os.Getenv("POSTGRES_URL"),
		BotToken:    os.Getenv("TELEGRAM_BOT_TOKEN"),
		WeatherKey:  os.Getenv("OPENWEATHER_API_KEY"),
	}
}
