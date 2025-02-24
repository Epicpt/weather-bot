package handlers

import (
	"fmt"
	"time"
	"weather-bot/internal/cache"
	"weather-bot/internal/database"
	"weather-bot/internal/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	volgograd   = "Волгоград"
	volgogradID = "498603"
)

// UserState представляет текущее состояние пользователя в диалоге
type UserState string

const (
	StateNone                     UserState = "none"
	StateAwaitingCityConfirmation UserState = "awaiting_city_confirmation"
	StateAwaitingCityInput        UserState = "awaiting_city_input"
	StateAwaitingCitySelection    UserState = "awaiting_city_selection"
)

func Update(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *database.Database, redisClient *cache.Cache) {
	// Получаем данные пользователя из Redis
	user, err := redisClient.GetUserFromRedis(update.Message.From.ID)
	if err != nil {
		log.Warn().Err(err).Msg("Ошибка при получении данных пользователя из Redis, пробуем получить из БД")

		// Если Redis недоступен, пытаемся получить пользователя из БД
		user, err = db.GetUserFromDB(update.Message.From.ID)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при получении данных пользователя из БД")
			// Если пользователь не найден в БД, создаем нового
			user = cache.NewUser(update.Message.From.ID, update.Message.From.FirstName)
			user.State = string(StateNone)
		}
	}
	// Если пользователь новый, инициализируем его
	if user == nil {
		user = cache.NewUser(update.Message.From.ID, update.Message.From.FirstName)
		user.State = string(StateNone)
	}

	// Обрабатываем сообщение в зависимости от текущего состояния пользователя
	switch UserState(user.State) {
	case StateNone:
		handleDefaultState(bot, update, user, db, redisClient)
	case StateAwaitingCityConfirmation:
		handleCityConfirmation(bot, update, user, db, redisClient)
	case StateAwaitingCityInput:
		handleCityInput(bot, update, user, db, redisClient)
	case StateAwaitingCitySelection:
		handleCitySelection(bot, update, user, db, redisClient)
	default:
		handleUnknownState(bot, update, user)
	}

	// Сохраняем обновленные данные пользователя в Redis
	err = redisClient.SaveUserToRedis(user)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при сохранении данных пользователя в Redis")
	}

}

func handleDefaultState(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *cache.User, db *database.Database, redisClient *cache.Cache) {
	switch update.Message.Text {
	case "/start":
		user.State = string(StateAwaitingCityConfirmation)
		sendMessage(bot, update.Message.Chat.ID, "Ваш город Волгоград?", cityConfirmationKeyboard())
	case "Узнать погоду", "/weather":
		forecast, err := weather.Get(user.CityID, redisClient, db)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при получении погоды")
			sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при получении погоды. Попробуйте повторить позже.", mainMenu())
			return
		}
		today := time.Now().UTC().Format("2006-01-02")
		msg := weather.FormatDailyForecast(user.City, forecast.FullDay[today])
		sendMessage(bot, update.Message.Chat.ID, msg, mainMenu())
	case "/weather5":
		forecast, err := weather.Get(user.CityID, redisClient, db)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при получении погоды")
			sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при получении погоды. Попробуйте повторить позже.", mainMenu())
			return
		}
		msg := weather.FormatFiveDayForecast(user.City, forecast.ShortDays)
		sendMessage(bot, update.Message.Chat.ID, msg, mainMenu())
	case "/city":
		user.State = string(StateAwaitingCityInput)
		sendMessage(bot, update.Message.Chat.ID, "Введите название вашего города:", tgbotapi.NewRemoveKeyboard(true))
	case "/notifications":
		// user.State = string(StateAwaitingCityInput) // TODO: state для выбора времени, вдруг введет неправильного формата время.
		sendMessage(bot, update.Message.Chat.ID, "Введите время в формате: часы.минуты (например: 09.15)", tgbotapi.NewRemoveKeyboard(true))
	default:
		sendMessage(bot, update.Message.Chat.ID, "Я не понимаю такую команду, выберите из меню.", mainMenu())
	}
}

func makeCityKeyboard(cities []cache.City) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton
	for _, city := range cities {
		text := fmt.Sprintf("%s|%d", city.Name, city.ID)
		if city.Region != "" {
			text = fmt.Sprintf("%s|%d|(%s)", city.Name, city.ID, city.Region)
		}
		row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(text))
		keyboard = append(keyboard, row)
	}
	keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Ввести название города заново.")))
	return tgbotapi.NewReplyKeyboard(keyboard...)
}

func handleUnknownState(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *cache.User) {
	user.State = string(StateNone)
	sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка. Начнем сначала.", startMenu())
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка отправки сообщения")
	}
}

func mainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Узнать погоду")))
}

func startMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/start")))
}
