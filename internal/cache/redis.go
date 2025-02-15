package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	c *redis.Client
}

func NewCache(c *redis.Client) *Cache {
	return &Cache{c: c}
}

// Инициализация PostgreSQL
func Init(addr string, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass, // no password set
		DB:       0,    // use default DB
	})
	return rdb

}

func (c *Cache) SaveUserToRedis(u *User) error {
	// Проверка инициализации клиента Redis
	if c.c == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	redisKey := fmt.Sprintf("user:%d", u.TgID)
	log.Info().Msgf("Redis key: %s", redisKey)

	// Формируем данные для записи
	userData := map[string]interface{}{
		"name":    u.Name,
		"city":    u.City,
		"city_id": u.CityID,
		"state":   u.State,
	}

	if u.FederalSubject != nil {
		userData["federal_subject"] = *u.FederalSubject
	}

	// Сохраняем в Redis
	err := c.c.HSet(context.Background(), redisKey, userData).Err()
	if err != nil {
		log.Error().Err(err).Msgf("Ошибка записи в Redis: %v", err)
		return fmt.Errorf("ошибка записи в Redis: %w", err)
	}

	log.Info().Msgf("Пользователь сохранён в Redis: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s", u.TgID, u.Name, u.City, u.CityID, u.State)
	return nil
}

func (c *Cache) GetUserFromRedis(userId int64) (*User, error) {
	// Проверка инициализации клиента Redis
	if c.c == nil {
		return nil, fmt.Errorf("Redis клиент не инициализирован")
	}

	redisKey := fmt.Sprintf("user:%d", userId)

	userData, err := c.c.HGetAll(context.Background(), redisKey).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из Redis: %w", err)
	}

	// Если пользователь не найден, возвращаем nil
	if len(userData) == 0 {
		return nil, nil
	}

	// Создаем и заполняем структуру User
	user := &User{
		TgID:   userId,
		Name:   userData["name"],
		City:   userData["city"],
		CityID: userData["city_id"],
		State:  userData["state"],
	}

	if federalSubject, ok := userData["federation_subject"]; ok {
		user.FederalSubject = &federalSubject
	}
	log.Info().Msgf("Пользователь получен из Redis: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s", user.TgID, user.Name, user.City, user.CityID, user.State)

	return user, nil
}
