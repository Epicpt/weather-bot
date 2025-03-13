package handlers

import (
	"fmt"
	"strconv"
	"time"
	"weather-bot/internal/app/reply"
	"weather-bot/internal/app/services"
	"weather-bot/internal/app/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

// UserState представляет текущее состояние пользователя в диалоге
type UserState string

const (
	StateNone                  UserState = "none"
	StateAwaitingCityInput     UserState = "awaiting_city_input"
	StateAwaitingCitySelection UserState = "awaiting_city_selection"
	StateAwaitingTimeInput     UserState = "awaiting_time_input"
)

func processMessage(ctx *Context) {

	// Обрабатываем сообщение в зависимости от текущего состояния пользователя
	switch UserState(ctx.user.State) {
	case StateNone:
		handleDefaultState(ctx)
	case StateAwaitingCityInput:
		handleCityInput(ctx)
	case StateAwaitingCitySelection:
		handleCitySelection(ctx)
	case StateAwaitingTimeInput:
		handleTimeInput(ctx)
	default:
		handleUnknownState(ctx)
	}
}

func handleDefaultState(ctx *Context) {
	switch ctx.text {
	case "/start":
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, startMessage(), tgbotapi.NewRemoveKeyboard(true))
	case "Узнать погоду", "/weather":
		forecast, err := weather.Get(ctx.user.CityID)
		if err != nil {
			log.Error().Err(err).Int64("user", ctx.user.TgID).Str("cityID", ctx.user.City).Msg("Ошибка при получении погоды")
			reply.Send().Message(ctx.user.ChatID, errorGetWeatherMessage(), mainMenu())
			return
		}
		today := time.Now().UTC().Format("2006-01-02")
		msg := weather.FormatDailyForecast(ctx.user.City, forecast.FullDay[today])
		reply.Send().Message(ctx.user.ChatID, msg, mainMenu())

		if ctx.user.Sticker {
			sticker := weather.Sticker(forecast.FullDay[today])
			reply.Send().Sticker(ctx.user.ChatID, sticker)
		}
	case "/weather5":
		forecast, err := weather.Get(ctx.user.CityID)
		if err != nil {
			log.Error().Err(err).Int64("user", ctx.user.TgID).Str("cityID", ctx.user.City).Msg("Ошибка при получении погоды")
			reply.Send().Message(ctx.user.ChatID, errorGetWeatherMessage(), mainMenu())
			return
		}
		msg := weather.FormatFiveDayForecast(ctx.user.City, forecast.ShortDays)
		reply.Send().Message(ctx.user.ChatID, msg, mainMenu())
	case "/city":
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, enterNameCityMessage(), tgbotapi.NewRemoveKeyboard(true))
	case "/notifications":
		existingTime, err := services.Global().GetUserNotificationTime(ctx.user.TgID)
		if err != nil {
			log.Error().Err(err).Int64("user", ctx.user.TgID).Msg("Ошибка при получении уведомления")
			reply.Send().Message(ctx.user.ChatID, "😢 Уведомления сейчас не работают. Попробуйте повторить позже.", mainMenu())
			return
		}

		ctx.user.State = string(StateAwaitingTimeInput)
		if existingTime != "" {
			// Уведомление уже есть, предлагаем изменить или удалить
			existingTime, err := strconv.ParseInt(existingTime, 10, 64)
			if err != nil {
				log.Error().Err(err).Int64("existing time", existingTime).Msg("Ошибка парсинга existingTime")
			}
			msg := fmt.Sprintf("❔ Вы уже получаете уведомления в %s.\nХотите изменить или удалить?", time.Unix(existingTime, 0).Format("15:04"))
			reply.Send().Message(ctx.user.ChatID, msg, notificationMenu())
		} else {
			// Уведомления нет, запрашиваем новое время
			reply.Send().Message(ctx.user.ChatID, enterNotificationTimeMessage(), tgbotapi.NewRemoveKeyboard(true))
		}
	case "/stickers":
		if ctx.user.Sticker {
			ctx.user.Sticker = false
			reply.Send().Message(ctx.user.ChatID, "Стикеры выключены ❌", mainMenu())
		} else {
			ctx.user.Sticker = true
			reply.Send().Message(ctx.user.ChatID, "Стикеры включены ✅", mainMenu())
		}

	default:
		reply.Send().Message(ctx.user.ChatID, "🤷‍♀️ Я не понимаю такую команду, выберите из меню.", mainMenu())
	}
}

func handleUnknownState(ctx *Context) {
	ctx.user.State = string(StateNone)
	reply.Send().Message(ctx.user.ChatID, "🔄 Произошла ошибка. Начнем сначала.", startMenu())
}
