package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"weather-bot/internal/app/reply"
	"weather-bot/internal/app/search"
	"weather-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func handleCityInput(ctx *Context) {
	cities, err := search.SearchCity(ctx.text)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при поиске города")
		reply.Send().Message(ctx.user.ChatID, "Произошла ошибка при поиске города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	if len(cities) == 1 {
		city := cities[0]
		ctx.user.Update(city.Name, strconv.Itoa(city.ID), string(StateNone), ctx.user.Sticker, city.Region)

		reply.Send().Message(ctx.user.ChatID, fmt.Sprintf("Отлично! Город %s сохранен.", city.Name), mainMenu())
		return
	}

	if len(cities) > 1 {
		keyboard := makeCityKeyboard(cities)
		ctx.user.State = string(StateAwaitingCitySelection)
		reply.Send().Message(ctx.user.ChatID, "Найдено несколько городов. Пожалуйста, выберите нужный:", keyboard)
		return
	}

	reply.Send().Message(ctx.user.ChatID, "Город не найден. Попробуйте ввести еще раз:", tgbotapi.NewRemoveKeyboard(true))
}

func handleCitySelection(ctx *Context) {

	if ctx.text == "Ввести название города заново." {
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, "Введите название вашего города:", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	parts := strings.Split(ctx.text, "|")
	if len(parts) < 2 || len(parts) > 3 {
		log.Error().Msg("Неверный формат выбранного города")
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, "Произошла ошибка при выборе города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	cityName := parts[0]
	cityID, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при парсинге ID города")
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, "Произошла ошибка при выборе города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))

		return
	}

	cities, err := search.SearchCity(cityName)
	if err != nil || len(cities) == 0 {
		log.Error().Err(err).Msg("Ошибка при поиске выбранного города")
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, "Произошла ошибка при выборе города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
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
		log.Error().Msg("Выбранный город не найден")
		ctx.user.State = string(StateAwaitingCityInput)
		reply.Send().Message(ctx.user.ChatID, "Произошла ошибка при выборе города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))

		return
	}

	ctx.user.Update(selectedCity.Name, strconv.Itoa(selectedCity.ID), string(StateNone), ctx.user.Sticker, selectedCity.Region)

	reply.Send().Message(ctx.user.ChatID, fmt.Sprintf("Отлично! Город %s сохранен.", selectedCity.Name), mainMenu())
}
