package cache

import (
	"context"
	"fmt"
	"strconv"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
)

var _ storage.UserStorage = (*Cache)(nil)

func (c *Cache) SaveUser(u *models.User) error {
	redisKey := fmt.Sprintf("user:%d", u.TgID)

	// Формируем данные для записи
	userData := map[string]interface{}{
		"chat_id": u.ChatID,
		"name":    u.Name,
		"city":    u.City,
		"city_id": u.CityID,
		"region":  u.Region,
		"state":   u.State,
		"sticker": u.Sticker,
	}

	// Сохраняем в Redis
	err := c.client.HSet(context.Background(), redisKey, userData).Err()
	if err != nil {
		log.Error().Err(err).Int64("userID", u.TgID).Msgf("Ошибка записи в Redis: %v", err)
		return fmt.Errorf("ошибка записи в Redis: %w", err)
	}

	log.Info().Msgf("Пользователь сохранён в Redis: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s, sticker=%v",
		u.TgID, u.Name, u.City, u.CityID, u.State, u.Sticker)
	return nil
}

func (c *Cache) GetUser(userId int64) (*models.User, error) {
	redisKey := fmt.Sprintf("user:%d", userId)

	userData, err := c.client.HGetAll(context.Background(), redisKey).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из Redis: %w", err)
	}

	// Если пользователь не найден, возвращаем nil
	if len(userData) == 0 {
		return nil, nil
	}

	chatId, err := strconv.ParseInt(userData["chat_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования числа из Redis: %w", err)
	}

	stickerBool, err := strconv.ParseBool(userData["sticker"])
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования булевого значения из Redis: %w", err)
	}

	// Создаем и заполняем структуру User
	user := &models.User{
		TgID:    userId,
		ChatID:  chatId,
		Name:    userData["name"],
		City:    userData["city"],
		CityID:  userData["city_id"],
		Region:  userData["region"],
		State:   userData["state"],
		Sticker: stickerBool,
	}

	log.Info().Msgf("Пользователь получен из Redis: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s", user.TgID, user.Name, user.City, user.CityID, user.State)

	return user, nil
}
