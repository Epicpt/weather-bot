package weather

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
	"weather-bot/internal/app/monitoring"
	"weather-bot/internal/app/services"
	"weather-bot/internal/models"

	"github.com/briandowns/openweathermap"
	"github.com/rs/zerolog/log"
)

// Временные промежутки
var dayParts = map[string][]int{
	"morning": {6, 9, 12},
	"day":     {12, 15, 18},
	"evening": {18, 21},
	"night":   {0, 3, 6},
}

func Get(cityID string) (*models.ProcessedForecast, error) {
	cityId, err := strconv.Atoi(cityID)
	if err != nil {
		return nil, fmt.Errorf("Неверный формат ID города: %v", err)
	}
	// Проверяем кеш
	if forecast, err := services.Global().GetWeather(cityId); err == nil {
		monitoring.WeatherCacheHitsTotal.Inc()
		return forecast, nil
	}
	log.Warn().Msg("не удалось получить forecast из хранилищ")

	// Получаем прогноз из OpenWeather
	processedForecast, err := GetNewWeather(cityId)
	if err != nil {
		monitoring.WeatherAPIErrorsTotal.Inc()
		return nil, fmt.Errorf("Не удалось получить погоду из OpenWeather: %v", err)
	}

	log.Info().Str("cityID", cityID).Msg("Новая погода для города получена")
	monitoring.WeatherCacheMissesTotal.Inc()
	return processedForecast, nil

}

func Update(cityIDs []string) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for _, cityID := range cityIDs {
		cityId, err := strconv.Atoi(cityID)
		if err != nil {
			return fmt.Errorf("Неверный формат ID города: %v", err)
		}

		<-ticker.C
		_, err = GetNewWeather(cityId)
		if err != nil {
			monitoring.WeatherAPIErrorsTotal.Inc()
			return err
		}
	}
	return nil
}

func GetNewWeather(cityID int) (*models.ProcessedForecast, error) {
	monitoring.WeatherAPIRequestsTotal.Inc()
	// Запрашиваем прогноз из OpenWeather
	forecastData, err := fetchWeatherFromAPI(cityID)
	if err != nil {
		return nil, err
	}

	// Обрабатываем прогноз сразу на 5 дней
	processedForecast, err := processWeatherData(forecastData)
	if err != nil {
		return nil, err
	}

	// Сохраняем (на 25 часов)
	if err = services.Global().SaveWeather(cityID, processedForecast); err != nil {
		log.Error().Err(err).Int("cityID", cityID).Msg("Error saving weather")
	}

	return processedForecast, nil
}

func fetchWeatherFromAPI(cityID int) (*openweathermap.ForecastWeatherData, error) {
	// Инициализируем клиент OpenWeather
	owm, err := openweathermap.NewForecast("5", "C", "ru", os.Getenv("OPENWEATHER_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации OpenWeather API: %w", err)
	}

	// Запрашиваем прогноз для города
	err = owm.DailyByID(cityID, 60) // 5-дневный прогноз
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса погоды: %w", err)
	}

	return owm, nil
}

func processWeatherData(forecast *openweathermap.ForecastWeatherData) (*models.ProcessedForecast, error) {
	forecastData, ok := forecast.ForecastWeatherJson.(*openweathermap.Forecast5WeatherData)
	if !ok {
		return nil, fmt.Errorf("не удалось преобразовать ForecastWeatherJson в Forecast5WeatherData")
	}

	// Создаём пустые карты для хранения прогноза
	fullDayForecasts := make(map[string]models.FullDayForecast)
	var shortDayForecasts []models.ShortDayForecast

	// Разбиваем прогноз по дням
	daysData := make(map[string][]openweathermap.Forecast5WeatherList)
	var dates []string

	for _, item := range forecastData.List {
		itemDate := time.Unix(int64(item.Dt), 0).UTC().Format("2006-01-02")
		daysData[itemDate] = append(daysData[itemDate], item)
	}

	// Сортируем даты по порядку
	for date := range daysData {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Обрабатываем каждый день
	for i, date := range dates {
		// Полный прогноз на 1 день (без ночи)
		dayForecast := processFullDayForecast(daysData[date])

		// Если есть следующий день — берём ночь оттуда
		if i+1 < len(dates) {
			nextDate := dates[i+1]
			nightForecast := calculateSummary(daysData[nextDate], dayParts["night"])
			dayForecast.Night = nightForecast
		}

		fullDayForecasts[date] = dayForecast

		// Краткий прогноз на 5 дней
		shortDayForecasts = append(shortDayForecasts, processShortDayForecast(date, daysData[date]))
	}

	return &models.ProcessedForecast{
		FullDay:   fullDayForecasts,
		ShortDays: shortDayForecasts,
	}, nil
}

func processFullDayForecast(data []openweathermap.Forecast5WeatherList) models.FullDayForecast {

	return models.FullDayForecast{
		Morning: calculateSummary(data, dayParts["morning"]),
		Day:     calculateSummary(data, dayParts["day"]),
		Evening: calculateSummary(data, dayParts["evening"]),
	}
}

func processShortDayForecast(date string, data []openweathermap.Forecast5WeatherList) models.ShortDayForecast {
	var tempSum float64
	var count int
	weatherCount := make(map[int]int)

	for _, item := range data {
		tempSum += item.Main.Temp
		count++

		// Подсчёт доминирующей погоды
		weatherCondition := item.Weather[0].ID
		weatherCount[weatherCondition]++
	}

	// Выбираем самую частую погоду
	dominantCondition, idCondition := getDominantCondition(weatherCount)

	if count == 0 {
		return models.ShortDayForecast{} // Если данных нет, возвращаем пустую структуру
	}

	return models.ShortDayForecast{
		Date:        date,
		Temperature: tempSum / float64(count),
		Condition:   dominantCondition,
		ConditionId: idCondition,
	}
}
