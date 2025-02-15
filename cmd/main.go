package main

import (
	"fmt"
	"os"

	"weather-bot/internal/cache"
	"weather-bot/internal/database"
	"weather-bot/internal/handlers"
	"weather-bot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Инициализация логгера
	log := logger.InitLogger()
	log.Info().Msg("Logger initialized")

	err := godotenv.Load("../config/.env")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading.env file")
	}

	// Подключаемся к БД
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, dbPort, dbName)
	postgres, err := database.Init(dsn)
	defer postgres.Close()
	if err != nil {
		log.Fatal().Err(err).Msgf("Ошибка подключения к БД: %v", err)
	}
	log.Info().Msg("Connected to PostgreSQL")

	// Инициализация Redis
	redisdb := cache.Init(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"))
	defer redisdb.Close()
	log.Info().Msg("Connected to Redis")

	db := database.NewDatabase(postgres)
	rdb := cache.NewCache(redisdb)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal().Err(err).Msg("Telegram bot token not provided")
	}

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating bot")
	}

	log.Info().Msgf("Бот %s запущен", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Пропускаем неполные сообщения
			continue
		}

		handlers.Update(bot, update, db, rdb)
	}

}
