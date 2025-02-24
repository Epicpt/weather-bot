package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"weather-bot/internal/cache"
	"weather-bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func handleCityConfirmation(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *cache.User, db *database.Database, redisClient *cache.Cache) {
	switch update.Message.Text {
	case "Да":
		user.City = volgograd
		user.CityID = volgogradID
		user.State = string(StateNone)
		err := db.SaveUserToDB(user)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при сохранении пользователя в базе")
			sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при сохранении города. Попробуйте еще раз.", mainMenu())
			return
		}
		sendMessage(bot, update.Message.Chat.ID, "Отлично! Город Волгоград сохранен.", mainMenu())
	case "Нет":
		user.State = string(StateAwaitingCityInput)
		sendMessage(bot, update.Message.Chat.ID, "Введите название вашего города:", tgbotapi.NewRemoveKeyboard(true))
	default:
		sendMessage(bot, update.Message.Chat.ID, "Пожалуйста, выберите 'Да' или 'Нет'.", cityConfirmationKeyboard())
	}
}

func handleCityInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *cache.User, db *database.Database, redisClient *cache.Cache) {
	cityName := update.Message.Text
	cities, err := redisClient.FindCity(cityName)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при поиске города в Redis")
		sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при поиске города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	if len(*cities) == 1 {
		city := (*cities)[0]
		user.City = city.Name
		user.CityID = strconv.Itoa(city.ID)
		if city.Region != "" {
			user.Region = &city.Region
		}

		err = db.SaveUserToDB(user)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при сохранении пользователя в базе")
			sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при сохранении города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
			return
		}

		user.State = string(StateNone)
		sendMessage(bot, update.Message.Chat.ID, fmt.Sprintf("Отлично! Город %s сохранен.", city.Name), mainMenu())
		return
	}

	if len(*cities) > 1 {
		keyboard := makeCityKeyboard(*cities)
		user.State = string(StateAwaitingCitySelection)
		sendMessage(bot, update.Message.Chat.ID, "Найдено несколько городов. Пожалуйста, выберите нужный:", keyboard)
		return
	}

	sendMessage(bot, update.Message.Chat.ID, "Город не найден. Попробуйте ввести еще раз:", tgbotapi.NewRemoveKeyboard(true))
}

func handleCitySelection(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *cache.User, db *database.Database, redisClient *cache.Cache) {
	selectedCityText := update.Message.Text

	if selectedCityText == "Ввести название города заново." {
		user.State = string(StateAwaitingCityInput)
		sendMessage(bot, update.Message.Chat.ID, "Введите название вашего города:", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	parts := strings.Split(selectedCityText, "|")
	if len(parts) < 2 || len(parts) > 3 {
		log.Error().Msg("Неверный формат выбранного города")
		sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при выборе города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
		user.State = string(StateAwaitingCityInput)
		return
	}

	cityName := parts[0]
	cityID, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при парсинге ID города")
		sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при выборе города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
		user.State = string(StateAwaitingCityInput)
		return
	}

	cities, err := redisClient.FindCity(cityName)
	if err != nil || len(*cities) == 0 {
		log.Error().Err(err).Msg("Ошибка при поиске выбранного города")
		sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при выборе города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
		user.State = string(StateAwaitingCityInput)
		return
	}

	var selectedCity *cache.City
	for _, city := range *cities {
		if city.ID == cityID {
			selectedCity = &city
			break
		}
	}

	if selectedCity == nil {
		log.Error().Msg("Выбранный город не найден")
		sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при выборе города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
		user.State = string(StateAwaitingCityInput)
		return
	}

	user.City = selectedCity.Name
	user.CityID = strconv.Itoa(selectedCity.ID)
	if selectedCity.Region != "" {
		user.Region = &selectedCity.Region
	}

	user.State = string(StateNone)

	err = db.SaveUserToDB(user)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при сохранении пользователя в базе")
		sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при сохранении города. Попробуйте еще раз.", tgbotapi.NewRemoveKeyboard(true))
		return
	}

	sendMessage(bot, update.Message.Chat.ID, fmt.Sprintf("Отлично! Город %s сохранен.", selectedCity.Name), mainMenu())
}

func cityConfirmationKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Да"),
			tgbotapi.NewKeyboardButton("Нет"),
		),
	)
}
