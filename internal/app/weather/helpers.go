package weather

import (
	"time"
	"weather-bot/internal/models"

	"github.com/briandowns/openweathermap"
)

var weatherMapping = map[int]string{
	200: "Гроза с небольшим дождём", 201: "Гроза с дождём", 202: "Гроза с сильным дождём",
	210: "Небольшая гроза", 211: "Гроза", 212: "Сильная гроза", 221: "Местами гроза", 230: "Гроза с лёгкой моросью",
	231: "Гроза моросью", 232: "Гроза с сильной моросью",

	300: "Лёгкая морось", 301: "Морось", 302: "Сильная морось",
	310: "Лёгкий моросящий дождь", 311: "Моросящий дождь", 312: "Сильный моросящий дождь",
	313: "Кратковременный дождь с моросью", 314: "Кратковременный ливень с моросью", 321: "Кратковременная морось",

	500: "Лёгкий дождь", 501: "Умеренный дождь", 502: "Сильный дождь",
	503: "Очень сильный дождь", 504: "Ливень", 511: "Ледяной дождь",
	520: "Кратковременный лёгкий дождь", 521: "Кратковременный дождь", 522: "Кратковременный ливень",
	531: "Местами кратковременный дождь",

	600: "Небольшой снег", 601: "Снег", 602: "Сильный снег",
	611: "Мокрый снег", 612: "Лёгкий мокрый снег", 613: "Сильный мокрый снег",
	615: "Лёгкий дождь со снегом", 616: "Дождь со снегом", 620: "Кратковременный небольшой снег",
	621: "Кратковременный снег", 622: "Кратковременный сильный снег",

	701: "Туман", 711: "Дымка", 721: "Лёгкий туман", 731: "Песчаная буря",
	741: "Сильный туман", 751: "Песок", 761: "Пыль", 762: "Вулканический пепел",
	771: "Шквал", 781: "Торнадо",

	800: "Ясно",
	801: "Небольшая облачность", 802: "Переменная облачность", 803: "Облачно с прояснениями", 804: "Пасмурно",
}

var weatherPriority = map[int]int{
	200: 6, 201: 6, 202: 6, 210: 6, 211: 6, 212: 6, 221: 6, 230: 6, 231: 6, 232: 6, // Гроза ⚡
	300: 5, 301: 5, 302: 5, 310: 5, 311: 5, 312: 5, 313: 5, 314: 5, 321: 5, // Морось 🌫
	500: 4, 501: 4, 502: 4, 503: 4, 504: 4, 511: 4, 520: 4, 521: 4, 522: 4, 531: 4, // Дождь 🌧
	600: 3, 601: 3, 602: 3, 611: 3, 612: 3, 613: 3, 615: 3, 616: 3, 620: 3, 621: 3, 622: 3, // Снег ❄
	701: 2, 711: 2, 721: 2, 731: 2, 741: 2, 751: 2, 761: 2, 762: 2, 771: 2, 781: 2, // Атмосферные явления 🌪
	800: 1, 801: 2, 802: 2, 803: 2, 804: 2, // Облачность ☁
}

func getDominantCondition(weatherList map[int]int) (string, int) {

	var dominantID int
	maxCount := 0
	maxPriority := 0

	for id, cnt := range weatherList {
		priority := weatherPriority[id]

		if cnt > maxCount || (cnt == maxCount && priority > maxPriority) {
			dominantID = id
			maxCount = cnt
			maxPriority = priority
		}
	}
	// Возвращаем описание на русском
	if description, exists := weatherMapping[dominantID]; exists {
		return description, dominantID
	}

	return "Неизвестная погода", 0
}

// Функция для вычисления средних значений
func calculateSummary(data []openweathermap.Forecast5WeatherList, hours []int) models.WeatherSummary {
	var tempSum, feelsLikeSum, windSum float64
	var count int
	weatherCount := make(map[int]int)

	for _, item := range data {
		hour := time.Unix(int64(item.Dt), 0).UTC().Hour()
		if contains(hours, hour) {
			tempSum += item.Main.Temp
			feelsLikeSum += item.Main.FeelsLike
			windSum += item.Wind.Speed
			count++

			// Подсчёт доминирующей погоды
			weatherCondition := item.Weather[0].ID
			weatherCount[weatherCondition]++
		}
	}

	// Выбираем самую частую погоду
	dominantCondition, idCondition := getDominantCondition(weatherCount)

	if count == 0 {
		return models.WeatherSummary{} // Если данных нет, возвращаем пустую структуру
	}

	return models.WeatherSummary{
		Temperature: tempSum / float64(count),
		FeelsLike:   feelsLikeSum / float64(count),
		WindSpeed:   windSum / float64(count),
		Condition:   dominantCondition,
		ConditionId: idCondition,
	}
}

// Вспомогательная функция для проверки наличия элемента в слайсе
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
