package handlers

import "fmt"

func startMessage() string {
	return `👋 Привет! Я бот, который поможет вам быть в курсе погоды каждый день!

Я умею:
- Присылать прогноз погоды в выбранное вами время ⏰
- Показывать погоду прямо сейчас по запросу 🔍
- Работать с базой из более чем 1500 городов России и не только 🗺️

✏ Введите название вашего города, чтобы я мог отправлять актуальную погоду:`
}

func errorFindCityMessage() string {
	return "⛔️ Произошла ошибка при поиске города. Попробуйте еще раз:"
}

func successSaveCityMessage(name string) string {
	return fmt.Sprintf("🎉 Отлично! Город %s сохранен.", name)
}
func enterNotificationTimeMessage() string {
	return "✏ Введите время в формате: часы:минуты (например: 09:15)"
}
func enterNameCityMessage() string {
	return "✏ Введите название вашего города:"
}
func enterNameDiffCityMessage() string {
	return "✏ Введите название другого города (ваш город не изменится):"
}
func errorGetWeatherMessage() string {
	return "⛔️ Произошла ошибка при получении погоды. Попробуйте повторить позже."
}
