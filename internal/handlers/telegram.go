package handlers

import (
	"weather-bot/internal/cache"
	"weather-bot/internal/database"
	"weather-bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	volgograd   = "Волгоград"
	volgogradID = "472757"
)

// UserState представляет текущее состояние пользователя в диалоге
type UserState string

const (
	StateNone                     UserState = "none"
	StateAwaitingCityConfirmation UserState = "awaiting_city_confirmation"
	StateAwaitingCityInput        UserState = "awaiting_city_input"
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
	default:
		handleUnknownState(bot, update, user, redisClient)
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
	case "Узнать погоду":
		// Здесь будет логика получения погоды
		sendMessage(bot, update.Message.Chat.ID, "Погода пока не реализована, но скоро будет!", mainMenu())
	case "Поменять город":
		user.State = string(StateAwaitingCityInput)
		sendMessage(bot, update.Message.Chat.ID, "Введите название вашего города:", tgbotapi.NewRemoveKeyboard(true))
	default:
		sendMessage(bot, update.Message.Chat.ID, "Я не понимаю такую команду, выберите из меню.", mainMenu())
	}
}

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
		sendMessage(bot, update.Message.Chat.ID, "Введите название вашего города:", tgbotapi.NewRemoveKeyboard(true)) // проверка на вечное цукиеми
	default:
		sendMessage(bot, update.Message.Chat.ID, "Пожалуйста, выберите 'Да' или 'Нет'.", cityConfirmationKeyboard())
	}
}

func handleCityInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *cache.User, db *database.Database, redisClient *cache.Cache) {
	cityName := update.Message.Text
	// Здесь будет логика поиска города
	// Если город найден, сохраняем его и переходим в StateNone
	// Если город не найден, просим ввести снова
	// Пример:
	city, err := services.FindCityInOpenWeather(cityName)
	if err == nil {
		user.City = city.Name
		//user.CityID = city.ID
		user.State = string(StateNone)
		err := db.SaveUserToDB(user)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при сохранении пользователя в базе")
			sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка при сохранении города. Попробуйте еще раз.", mainMenu())
			return
		}
		sendMessage(bot, update.Message.Chat.ID, "Город успешно сохранен!", mainMenu())
	} else {
		sendMessage(bot, update.Message.Chat.ID, "Город не найден. Попробуйте ввести еще раз:", tgbotapi.NewRemoveKeyboard(true)) // выдача похожих через levenshtein
	}
}

func handleUnknownState(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *cache.User, redisClient *cache.Cache) { // перенаправлять на /start
	user.State = string(StateNone)
	sendMessage(bot, update.Message.Chat.ID, "Произошла ошибка. Начнем сначала.", mainMenu())
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка отправки сообщения")
	}
}

func cityConfirmationKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Да"),
			tgbotapi.NewKeyboardButton("Нет"),
		),
	)
}

func mainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Узнать погоду")))
}
