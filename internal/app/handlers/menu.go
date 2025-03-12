package handlers

import (
	"fmt"
	"weather-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func mainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Узнать погоду")))
}

func startMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/start")))
}

func makeCityKeyboard(cities []models.City) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton
	for _, city := range cities {
		text := fmt.Sprintf("%s|%d", city.Name, city.ID)
		if city.Region != "" {
			text = fmt.Sprintf("%s|%d|(%s)", city.Name, city.ID, city.Region)
		}
		row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(text))
		keyboard = append(keyboard, row)
	}
	keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("🔄 Ввести название города заново.")))
	return tgbotapi.NewReplyKeyboard(keyboard...)
}

func notificationMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✏ Изменить"),
			tgbotapi.NewKeyboardButton("❌ Удалить"),
			tgbotapi.NewKeyboardButton("↩ Отмена"),
		),
	)
}
