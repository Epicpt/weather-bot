package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FindCity ищет город в redis
func (c *Cache) FindCity(cityName string) (*[]City, error) {
	caser := cases.Title(language.Russian)
	cityName = caser.String(strings.ToLower(cityName))

	cities, err := c.GetCities(cityName)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения города из resis: %w", err)
	}

	if cities == nil || len(*cities) == 0 {
		closestMatch, err := c.FindTop3ClosestCities(cityName)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения похожих городов: %w", err)
		}
		if closestMatch == nil {
			return nil, fmt.Errorf("Похожие города с таким названием не найдены")
		}
		return closestMatch, nil
	}

	return cities, nil

}

// FindTop3ClosestCities находит 3 похожих города
func (c *Cache) FindTop3ClosestCities(input string) (*[]City, error) {
	if input == "" {
		return nil, fmt.Errorf("ввод не должен быть пустым")
	}

	type cityDistance struct {
		city     string
		distance int
	}

	var distances []cityDistance

	cities, err := c.GetAllCitiesKeys()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ключей городов: %w", err)
	}

	// Считаем расстояния
	for _, city := range cities {
		distance := levenshtein.DistanceForStrings([]rune(input), []rune(city), levenshtein.DefaultOptions)
		distances = append(distances, cityDistance{city: city, distance: distance})
	}

	// Сортируем по возрастанию расстояния
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	var closestCities []City

	for i := 0; i < 3 && i < len(distances); i++ {
		log.Debug().Msgf("Поиск города: %s", distances[i].city)
		citiesClose, err := c.GetCities(distances[i].city)
		if err != nil {
			log.Error().Err(err).Msgf("Ошибка при получении города %s из Redis", distances[i].city)
			continue // Пропускаем ошибку, но продолжаем
		}

		if citiesClose == nil || len(*citiesClose) == 0 {
			log.Warn().Msgf("Город %s не найден в Redis", distances[i].city)
			continue
		}
		log.Debug().Msgf("Найдено %d городов для %s", len(*citiesClose), distances[i].city)

		for _, c := range *citiesClose {
			closestCities = append(closestCities, c)
		}

	}

	if len(closestCities) == 0 {
		return nil, fmt.Errorf("похожие города не найдены")
	}

	return &closestCities, nil
}

func (c *Cache) LoadCities(jsonFile string) error {
	file, err := os.Open(jsonFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var cities []City
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cities)
	if err != nil {
		return err
	}

	for _, city := range cities {
		err = c.saveCityToRedis(city)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) saveCityToRedis(city City) error {
	if c.c == nil {
		return fmt.Errorf("Redis клиент не инициализирован")
	}

	redisKey := fmt.Sprintf("city:%s", city.Name)

	// Преобразуем структуру в JSON
	cityData, err := json.Marshal(city)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при сериализации города")
		return fmt.Errorf("ошибка при сериализации города: %w", err)
	}

	// Проверяем, существует ли уже город с таким ID
	existingCities, err := c.c.LRange(context.Background(), redisKey, 0, -1).Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при проверке существующих городов")
		return fmt.Errorf("ошибка при проверке существующих городов: %w", err)
	}

	for i, existingCityData := range existingCities {
		var existingCity City
		if err := json.Unmarshal([]byte(existingCityData), &existingCity); err != nil {
			continue
		}
		if existingCity.ID == city.ID {
			// Город с таким ID уже существует, обновляем его
			err = c.c.LSet(context.Background(), redisKey, int64(i), string(cityData)).Err()
			if err != nil {
				log.Error().Err(err).Msg("Ошибка при обновлении города в Redis")
				return fmt.Errorf("ошибка при обновлении города в Redis: %w", err)
			}
			return nil
		}
	}

	// Если город не найден, добавляем новый
	err = c.c.RPush(context.Background(), redisKey, cityData).Err()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка записи в Redis")
		return fmt.Errorf("ошибка записи в Redis: %w", err)
	}

	return nil
}

func (c *Cache) GetCities(city string) (*[]City, error) {
	// Проверка инициализации клиента Redis
	if c.c == nil {
		return nil, fmt.Errorf("Redis клиент не инициализирован")
	}

	seen := make(map[string]bool)
	var result []City
	redisKey := fmt.Sprintf("city:%s", city)

	citiesData, err := c.c.LRange(context.Background(), redisKey, 0, -1).Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка получения городов из Redis")
		return nil, fmt.Errorf("ошибка получения городов из Redis: %w", err)
	}
	for _, data := range citiesData {
		var city City
		err := json.Unmarshal([]byte(data), &city)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка десериализации города")
			continue
		}
		if seen[city.Region] {
			continue
		}

		seen[city.Region] = true
		result = append(result, city)
	}
	return &result, nil
}

func (c *Cache) GetAllCitiesKeys() ([]string, error) {
	if c.c == nil {
		return nil, fmt.Errorf("Redis клиент не инициализирован")
	}

	// Получаем все ключи, соответствующие шаблону "city:*"
	keys, err := c.c.Keys(context.Background(), "city:*").Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ключей из Redis: %w", err)
	}

	// Удаляем префикс "city:" из каждого ключа
	cities := make([]string, 0, len(keys))
	for _, key := range keys {
		cityName := strings.TrimPrefix(key, "city:")
		cities = append(cities, cityName)
	}

	return cities, nil
}

func (c *Cache) GetAllCitiesIds() ([]string, error) {
	ctx := context.Background()

	// Получаем все ключи, соответствующие шаблону "user:*"
	userKeys, err := c.c.Keys(ctx, "user:*").Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка получения user-ключей из Redis")
		return nil, err
	}

	var cityIDs []string
	seenCities := make(map[string]bool)

	// Извлекаем `city_id` у каждого пользователя
	for _, key := range userKeys {
		cityID, err := c.c.HGet(ctx, key, "city_id").Result()
		if err == nil && !seenCities[cityID] {
			cityIDs = append(cityIDs, cityID)
			seenCities[cityID] = true
		}
	}

	if len(cityIDs) == 0 {
		return nil, fmt.Errorf("нет данных в Redis, нужно запросить из БД")
	}

	return cityIDs, nil
}
