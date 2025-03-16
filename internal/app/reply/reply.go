package reply

import (
	"time"
	"weather-bot/internal/app/weather"
	"weather-bot/internal/models"
)

type Sender interface {
	Message(chatID int64, text string, keyboard any)
	Sticker(chatID int64, stickerID string)
}

var sender Sender

func Init(s Sender) {
	sender = s
}

func Send() Sender {
	return sender
}

func SendDailyWeather(user *models.User, forecast *models.ProcessedForecast) {
	today := time.Now().UTC().Format("2006-01-02")

	msg := weather.FormatDailyForecast(user.City, forecast.FullDay[today])
	Send().Message(user.ChatID, msg, nil)

	if user.Sticker {
		sticker := weather.Sticker(forecast.FullDay[today])
		Send().Sticker(user.ChatID, sticker)
	}
}
