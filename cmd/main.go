package main

import (
	"os"
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

		log.Info().Msgf("Получено сообщение от %s: %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, "+update.Message.From.FirstName+"! Обновляю функционал бота!")
		bot.Send(msg)
		// Отправляем ответный текст в чат
		//...
		// Вы можете добавлять другие действия в зависимости от входящего сообщения
		//...
		// Обработка ошибок
		//...
		// Остановка бота при получении конфигурации /stop
	}

}
