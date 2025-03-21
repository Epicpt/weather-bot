package weather

import (
	"fmt"
	"time"
	"weather-bot/internal/models"
)

func FormatDailyForecast(city string, forecast models.FullDayForecast) string {
	message := fmt.Sprintf("🌤 <b>Прогноз на сегодня (%s):</b>\n", city)

	if (forecast.Morning != models.WeatherSummary{}) {
		message += fmt.Sprintf("🥱 <b>Утро:</b> %.f°C, ощущается как %.1f°C, %s %s\n",
			forecast.Morning.Temperature, forecast.Morning.FeelsLike, forecast.Morning.Condition, getWeatherEmoji(forecast.Morning.ConditionId))
	}
	if (forecast.Day != models.WeatherSummary{}) {
		message += fmt.Sprintf("🌞 <b>День:</b> %.f°C, ощущается как %.1f°C, %s %s\n",
			forecast.Day.Temperature, forecast.Day.FeelsLike, forecast.Day.Condition, getWeatherEmoji(forecast.Day.ConditionId))
	}

	if (forecast.Evening != models.WeatherSummary{}) {
		message += fmt.Sprintf("🌚 <b>Вечер:</b> %.f°C, ощущается как %.1f°C, %s %s\n",
			forecast.Evening.Temperature, forecast.Evening.FeelsLike, forecast.Evening.Condition, getWeatherEmoji(forecast.Evening.ConditionId))
	}

	if (forecast.Night != models.WeatherSummary{}) {
		message += fmt.Sprintf("🌙 <b>Ночь:</b> %.f°C, ощущается как %.1f°C, %s %s",
			forecast.Night.Temperature, forecast.Night.FeelsLike, forecast.Night.Condition, getWeatherEmoji(forecast.Night.ConditionId))
	}

	return message
}

func getWeatherEmoji(conditionId int) string {
	if conditionId == 800 {
		return "☀️"
	}
	id := conditionId / 100

	switch id {
	case 2:
		return "⛈"
	case 3:
		return "🌦"
	case 5:
		return "🌧"
	case 6:
		return "❄️"
	case 7:
		return "🌫"
	case 8:
		return "☁️"
	default:
		return ""
	}
}

func FormatFiveDayForecast(city string, forecasts []models.ShortDayForecast) string {
	message := fmt.Sprintf("🌤 <b>Прогноз на 5 дней (%s):</b>\n", city)

	for _, f := range forecasts {
		emoji := getWeatherEmoji(f.ConditionId)
		message += fmt.Sprintf("🗓 <b>%s:</b> %.f°C, %s %s\n",
			formatDate(f.Date), f.Temperature, f.Condition, emoji)
	}

	return message
}

func formatDate(dateStr string) string {
	days := map[string]string{"Monday": "Пн", "Tuesday": "Вт", "Wednesday": "Ср", "Thursday": "Чт", "Friday": "Пт", "Saturday": "Сб", "Sunday": "Вс"}
	t, _ := time.Parse("2006-01-02", dateStr)
	return days[t.Weekday().String()]
}
