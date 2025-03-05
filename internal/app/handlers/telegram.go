package handlers

import (
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

func Update(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Получаем данные пользователя из хранилища
	userService := services.Global()
	user, err := userService.GetUser(update.Message.From.ID)
	if err != nil {
		log.Warn().Err(err).Msg("Ошибка при получении данных пользователя из хранилища")
	}

	// Если пользователь новый, инициализируем его
	if user == nil {
		log.Info().Msgf("Новый пользователь %s!", update.Message.From.FirstName)
		user = models.NewUser(update.Message.From.ID, update.Message.Chat.ID, update.Message.From.FirstName, string(StateNone))
	}

	ctx := &Context{
		bot:  bot,
		user: user,
		text: update.Message.Text,
	}

	processMessage(ctx)

	// Сохраняем обновленные данные пользователя
	if err = userService.SaveUser(user); err != nil {
		log.Error().Err(err).Msg("Ошибка при сохранении пользователя в хранилище")
	}

}
