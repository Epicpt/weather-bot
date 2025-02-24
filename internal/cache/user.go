package cache

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type User struct {
	TgID             int64   `json:"tg_id"`
	Name             string  `json:"name"`
	City             string  `json:"city"`
	CityID           string  `json:"city_id"`
	Region           *string `json:"federal_subject,omitempty"`
	NotificationTime *string `json:"notification_time,omitempty"`
	State            string  `json:"state"`
}

func NewUser(tgID int64, name string) *User {
	return &User{
		TgID: tgID,
		Name: name,
	}
}

func (c *Cache) SaveUserToRedis(u *User) error {
	// Проверка инициализации клиента Redis
	if c.c == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	redisKey := fmt.Sprintf("user:%d", u.TgID)

	// Формируем данные для записи
	userData := map[string]interface{}{
		"name":    u.Name,
		"city":    u.City,
		"city_id": u.CityID,
		"state":   u.State,
	}

	if u.Region != nil {
		userData["region"] = *u.Region
	}

	if u.NotificationTime != nil {
		userData["notification_time"] = *u.NotificationTime
	}

	// Сохраняем в Redis
	err := c.c.HSet(context.Background(), redisKey, userData).Err()
	if err != nil {
		log.Error().Err(err).Msgf("Ошибка записи в Redis: %v", err)
		return fmt.Errorf("ошибка записи в Redis: %w", err)
	}

	log.Info().Msgf("Пользователь сохранён в Redis: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s, notification_time=%v",
		u.TgID, u.Name, u.City, u.CityID, u.State, u.NotificationTime)
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
		user.Region = &federalSubject
	}
	if notificationTime, ok := userData["notification_time"]; ok {
		user.NotificationTime = &notificationTime
	}
	log.Info().Msgf("Пользователь получен из Redis: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s", user.TgID, user.Name, user.City, user.CityID, user.State)

	return user, nil
}
