package reply

import (
	"strings"
	"time"
	"weather-bot/internal/app/services"
	"weather-bot/internal/app/weather"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
)

type Sender interface {
	Message(chatID int64, text string, keyboard any) error
	Sticker(chatID int64, stickerID string) error
}

var sender Sender

func Init(s Sender) {
	sender = s
}

func Send() Sender {
	return sender
}

func SendDailyWeather(user *models.User, forecast *models.ProcessedForecast) error {
	today := time.Now().UTC().Format("2006-01-02")

	msg := weather.FormatDailyForecast(user.City, forecast.FullDay[today])
	err := Send().Message(user.ChatID, msg, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Forbidden: bot was blocked by the user") {
			log.Warn().Err(err).Msgf("reply - SendDailyWeather - Пользователь %d заблокировал бота", user.TgID)
			if err := services.Global().NotificationService.RemoveUserNotification(user.TgID); err != nil {
				log.Error().Err(err).Int64("user", user.TgID).Msg("reply - SendDailyWeather - Ошибка при удалении уведомления")
			}

		} else {
			log.Error().Err(err).Int64("user", user.TgID).Msg("reply - SendDailyWeather - Ошибка при отправке сообщения")
		}

		return err
	}

	if user.Sticker {
		sticker := weather.Sticker(forecast.FullDay[today])
		err := Send().Sticker(user.ChatID, sticker)
		if err != nil {
			log.Error().Err(err).Int64("user", user.TgID).Str("sticker", sticker).Msg("reply - SendDailyWeather - Ошибка при отправке стикера")
		}
	}
	return nil
}
