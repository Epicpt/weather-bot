package handlers

import (
	"weather-bot/internal/app/monitoring"
	"weather-bot/internal/app/services"
	"weather-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type Context struct {
	bot  *tgbotapi.BotAPI
	user *models.User
	text string
}

func Update(update tgbotapi.Update) {
	monitoring.BotRequestsTotal.Inc()
	monitoring.UpdateUniqueUsers(update.Message.From.ID)

	// Получаем данные пользователя из хранилища
	userService := services.Global()
	user, err := userService.GetUser(update.Message.From.ID)
	if err != nil {
		log.Warn().Err(err).Int64("id", update.Message.From.ID).Str("user", update.Message.From.FirstName).Msg("Ошибка при получении данных пользователя из хранилища")
	}

	// Если пользователь новый, инициализируем его
	if user == nil {
		user = models.NewUser(update.Message.From.ID, update.Message.Chat.ID, update.Message.From.FirstName, string(StateNone))
		log.Info().Int64("id", user.TgID).Msgf("Новый пользователь %s!", user.Name)
	}

	ctx := &Context{
		user: user,
		text: update.Message.Text,
	}

	log.Info().Int64("id", user.TgID).Str("user", user.Name).Str("username", update.Message.From.UserName).Str("city", user.City).Str("state", user.State).Bool("sticker", user.Sticker).
		Msgf("Пользователь отправил сообщение: %s", update.Message.Text)

	processMessage(ctx)

	// Сохраняем обновленные данные пользователя
	if err = userService.SaveUser(user); err != nil {
		monitoring.BotErrorsTotal.Inc()
		log.Error().Err(err).Int64("id", user.TgID).Msg("Ошибка при сохранении пользователя в хранилище")
	}

}
