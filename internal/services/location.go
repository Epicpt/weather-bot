package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

const openWeatherGeocodeURL = "http://api.openweathermap.org/geo/1.0/direct"

type City struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

// FindCityInOpenWeather ищет город в OpenWeather API
func FindCityInOpenWeather(cityName string) (*City, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("openweathermap api key not found")
	}
	url := fmt.Sprintf("%s?q=%s&limit=5&appid=%s&lang=ru", openWeatherGeocodeURL, cityName, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cities []City
	if err := json.Unmarshal(body, &cities); err != nil {
		return nil, err
	}

	if len(cities) == 0 {
		return nil, errors.New("город не найден")
	}

	return &cities[0], nil
}

// FindClosestMatch находит город с минимальной ошибкой
func FindClosestMatch(input string, cities []string) string {
	minDistance := 100
	closestMatch := ""

	for _, city := range cities {
		distance := levenshtein.DistanceForStrings([]rune(input), []rune(city), levenshtein.DefaultOptions)
		if distance < minDistance {
			minDistance = distance
			closestMatch = city
		}
	}

	if minDistance <= 2 {
		return closestMatch
	}

	return ""
}

// GetCitySuggestions возвращает список похожих городов
func GetCitySuggestions(input string) []string {
	allCities := []string{"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург", "Казань"} // json all
	var suggestions []string

	for _, city := range allCities {
		distance := levenshtein.DistanceForStrings([]rune(input), []rune(city), levenshtein.DefaultOptions)
		if distance <= 2 {
			suggestions = append(suggestions, city)
		}
	}

	return suggestions
}
