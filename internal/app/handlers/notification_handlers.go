package handlers

import (
	"fmt"
	"regexp"
	"time"
	"weather-bot/internal/app/jobs"
	"weather-bot/internal/app/reply"
	"weather-bot/internal/app/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func handleTimeInput(ctx *Context) {

	switch ctx.text {
	case "↩ Отмена":
		ctx.user.State = string(StateNone)
		reply.Send().Message(ctx.user.ChatID, "Отменено.", mainMenu())
		return
	case "❌ Удалить":
		ctx.user.State = string(StateNone)
		err := services.Global().RemoveUserNotification(ctx.user.TgID)
		if err != nil {
			log.Error().Err(err).Int64("user", ctx.user.TgID).Msg("Ошибка при удалении уведомления")
			reply.Send().Message(ctx.user.ChatID, "❌ Ошибка при удалении уведомления.", mainMenu())
		} else {
			reply.Send().Message(ctx.user.ChatID, "✅ Уведомление удалено.", mainMenu())
		}
		return
	case "✏ Изменить":
		reply.Send().Message(ctx.user.ChatID, enterNotificationTimeMessage(), tgbotapi.NewRemoveKeyboard(true))
		return
	}

	if isValidTime(ctx.text) {
		// Парсим `HH:MM` в `time.Time`
		notifTime, err := time.Parse("15:04", ctx.text)
		if err != nil {
			log.Error().Err(err).Int64("user", ctx.user.TgID).Time("notification time", notifTime).Msgf("Ошибка парсинга времени: %s", ctx.text)
		}
		err = jobs.ScheduleUserUpdate(ctx.user.TgID, notifTime)
		if err != nil {
			log.Error().Err(err).Int64("user", ctx.user.TgID).Msg("Ошибка при добавлении уведомлений")
		}
		log.Info().Int64("user", ctx.user.TgID).Msgf("Время %s для юзера %s сохранено", ctx.text, ctx.user.Name)
		ctx.user.State = string(StateNone)
		reply.Send().Message(ctx.user.ChatID, fmt.Sprintf("🎉 Отлично! Завтра в %s вам придет прогноз погоды.", ctx.text), mainMenu())

	} else {
		reply.Send().Message(ctx.user.ChatID, "⛔️ Неверный формат времени (часы:минуты). Попробуйте ввести еще раз.", tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("↩ Отмена"))))
	}

}

func isValidTime(input string) bool {
	// Проверяем формат через регулярку "HH:MM"
	matched, err := regexp.MatchString(`^([01]\d|2[0-3]):([0-5]\d)$`, input)
	if err != nil {
		log.Error().Err(err).Str("notification time", input).Msg("Ошибка валидации времени")
		return false
	}
	return matched
}
