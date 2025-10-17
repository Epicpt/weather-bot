package handlers

import (
	"regexp"
	"strconv"
	"strings"
	"weather-bot/internal/app/reply"
	"weather-bot/internal/app/search"
	"weather-bot/internal/app/weather"
	"weather-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func handleCityInput(ctx *Context) {
	if !IsValidCity(ctx.text) {
		reply.Send().Message(ctx.user.ChatID, "⛔️ Принимается название города только на кириллице. Попробуйте еще раз:", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	cities, err := search.SearchCity(ctx.text)
	if err != nil {
		log.Error().Err(err).Int64("user", ctx.user.TgID).Str("city", ctx.text).Msg("Ошибка при поиске города")
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), tgbotapi.NewRemoveKeyboard(true))
		return
	}

	if len(cities) == 1 {
		log.Info().Int64("user", ctx.user.TgID).Str("city", ctx.text).Msg("Пользователь выбрал город")
		city := cities[0]
		ctx.user.Update(city.Name, strconv.Itoa(city.ID), string(StateNone), ctx.user.Sticker, city.Region)

		reply.Send().Message(ctx.user.ChatID, successSaveCityMessage(city.Name), mainMenu())
		return
	}

	if len(cities) > 1 {
		keyboard := makeCityKeyboard(cities)
		ctx.user.State = string(StateAwaitingCitySelection)
		reply.Send().Message(ctx.user.ChatID, "🔍 Найдено несколько городов. Пожалуйста, выберите нужный:", keyboard)
		return
	}

	reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), tgbotapi.NewRemoveKeyboard(true))
}

func IsValidCity(city string) bool {
	r := regexp.MustCompile(`^[а-яА-ЯёЁ\s-]+$`)
	return r.MatchString(city)
}

func handleCitySelection(ctx *Context) {

	if ctx.text == "🔄 Ввести название города заново." {
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, enterNameCityMessage(), tgbotapi.NewRemoveKeyboard(true))
		return
	}

	parts := strings.Split(ctx.text, "|")
	if len(parts) < 2 || len(parts) > 3 {
		log.Error().Str("city", ctx.text).Msg("Неверный формат выбранного города")
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), tgbotapi.NewRemoveKeyboard(true))
		return
	}

	cityName := parts[0]
	cityID, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Error().Int("cityID", cityID).Err(err).Msg("Ошибка при парсинге ID города")
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), tgbotapi.NewRemoveKeyboard(true))

		return
	}

	cities, err := search.SearchCity(cityName)
	if err != nil || len(cities) == 0 {
		log.Error().Int("cityID", cityID).Str("city", cityName).Err(err).Msg("Ошибка при поиске выбранного города")
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), tgbotapi.NewRemoveKeyboard(true))
		return
	}

	var selectedCity *models.City
	for _, city := range cities {
		if city.ID == cityID {
			selectedCity = &city
			break
		}
	}

	if selectedCity == nil {
		log.Error().Int("cityID", cityID).Str("city", cityName).Msg("Выбранный город не найден")
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), tgbotapi.NewRemoveKeyboard(true))

		return
	}

	ctx.user.Update(selectedCity.Name, strconv.Itoa(selectedCity.ID), string(StateNone), ctx.user.Sticker, selectedCity.Region)

	log.Info().Int64("user", ctx.user.TgID).Str("city", selectedCity.Name).Msg("Пользователь выбрал город")

	reply.Send().Message(ctx.user.ChatID, successSaveCityMessage(selectedCity.Name), mainMenu())
}

func handleDiffCityInput(ctx *Context) {
	if ctx.text == "↩ Отмена" {
		ctx.user.State = string(StateNone)
		reply.Send().Message(ctx.user.ChatID, "Отменено.", mainMenu())
		return
	}
	if !IsValidCity(ctx.text) {
		reply.Send().Message(ctx.user.ChatID, "⛔️ Принимается название города только на кириллице. Попробуйте еще раз:", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	cities, err := search.SearchCity(ctx.text)
	if err != nil {
		log.Error().Err(err).Int64("user", ctx.user.TgID).Str("city", ctx.text).Msg("Ошибка при поиске города")
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), cancelMenu())
		return
	}

	if len(cities) == 1 {
		ctx.user.State = string(StateNone)
		city := cities[0]
		forecast, err := weather.Get(strconv.Itoa(city.ID))
		if err != nil {
			log.Error().Err(err).Int64("user", ctx.user.TgID).Int("cityID", city.ID).Msg("Ошибка при получении погоды")
			reply.Send().Message(ctx.user.ChatID, errorGetWeatherMessage(), mainMenu())
			return
		}
		msg := weather.FormatFiveDayForecast(city.Name, forecast.ShortDays)

		reply.Send().Message(ctx.user.ChatID, msg, mainMenu())
		return
	}

	if len(cities) > 1 {
		keyboard := makeCityKeyboard(cities)
		ctx.user.State = string(StateAwaitingDiffCitySelection)
		reply.Send().Message(ctx.user.ChatID, "🔍 Найдено несколько городов. Пожалуйста, выберите нужный:", keyboard)
		return
	}

	reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), cancelMenu())
}

func handleDiffCitySelection(ctx *Context) {

	if ctx.text == "🔄 Ввести название города заново." {
		ctx.user.State = string(StateAwaitingDiffCityInput)
		reply.Send().Message(ctx.user.ChatID, enterNameDiffCityMessage(), cancelMenu())
		return
	}

	parts := strings.Split(ctx.text, "|")
	if len(parts) < 2 || len(parts) > 3 {
		log.Error().Str("city", ctx.text).Msg("Неверный формат выбранного города")
		ctx.user.State = string(StateAwaitingDiffCityInput)
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), cancelMenu())
		return
	}

	cityName := parts[0]
	cityID, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Error().Int("cityID", cityID).Err(err).Msg("Ошибка при парсинге ID города")
		ctx.user.State = string(StateAwaitingDiffCityInput)
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), cancelMenu())

		return
	}

	cities, err := search.SearchCity(cityName)
	if err != nil || len(cities) == 0 {
		log.Error().Int("cityID", cityID).Str("city", cityName).Err(err).Msg("Ошибка при поиске выбранного города")
		ctx.user.State = string(StateAwaitingDiffCityInput)
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), cancelMenu())
		return
	}

	var selectedCity *models.City
	for _, city := range cities {
		if city.ID == cityID {
			selectedCity = &city
			break
		}
	}

	if selectedCity == nil {
		log.Error().Int("cityID", cityID).Str("city", cityName).Msg("Выбранный город не найден")
		ctx.user.State = string(StateAwaitingDiffCityInput)
		reply.Send().Message(ctx.user.ChatID, errorFindCityMessage(), cancelMenu())

		return
	}

	ctx.user.State = string(StateNone)
	forecast, err := weather.Get(strconv.Itoa(selectedCity.ID))
	if err != nil {
		log.Error().Err(err).Int64("user", ctx.user.TgID).Int("cityID", selectedCity.ID).Msg("Ошибка при получении погоды")
		reply.Send().Message(ctx.user.ChatID, errorGetWeatherMessage(), mainMenu())
		return
	}
	msg := weather.FormatFiveDayForecast(selectedCity.Name, forecast.ShortDays)

	reply.Send().Message(ctx.user.ChatID, msg, mainMenu())
}
