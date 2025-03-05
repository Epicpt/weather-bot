package cache

import (
	"context"
	"fmt"
	"strconv"
	"weather-bot/internal/app/storage"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var _ storage.NotificationStorage = (*Cache)(nil)

func (c *Cache) ScheduleWeatherUpdate(executeAt int64) error {
	ctx := context.Background()
	// ID задачи
	jobID := fmt.Sprintf("weather_update:%d", executeAt)

	// Добавляем задачу в Redis Stream
	_, err := c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "weather_updates",
		Values: map[string]interface{}{
			"job_id":    jobID,
			"executeAt": executeAt,
		},
	}).Result()

	if err != nil {
		return fmt.Errorf("Ошибка записи задачи на обновление погоды в Redis: %w", err)
	}
	return nil
}

func (c *Cache) RemoveWeatherUpdate() error {
	ctx := context.Background()

	messages, err := c.client.XRange(ctx, "weather_updates", "-", "+").Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения из Redis Stream")
		return err
	}

	for _, msg := range messages {
		_, err := c.client.XDel(ctx, "weather_updates", msg.ID).Result()
		if err != nil {
			log.Error().Err(err).Msgf("Ошибка удаления задачи %s", msg.ID)
			return err
		}
	}
	return nil
}

func (c *Cache) GetScheduleWeatherUpdate() ([]redis.XStream, error) {
	ctx := context.Background()

	streams, err := c.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"weather_updates", "0"},
		Count:   1,
		Block:   0, // Блокируемся и ждём новую задачу
	}).Result()

	if err != nil {
		return nil, err
	}
	return streams, nil
}

func (c *Cache) ScheduleUserNotification(userID int64, executeAt int64) error {
	ctx := context.Background()

	_, err := c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "user_notifications",
		Values: map[string]interface{}{
			"user_id":   userID,
			"executeAt": executeAt,
		},
	}).Result()

	if err != nil {
		return fmt.Errorf("Ошибка записи задачи уведомления в Redis Stream: %w", err)
	}
	return nil
}

func (c *Cache) RemoveUserNotification(userID int64) error {
	ctx := context.Background()

	messages, err := c.client.XRange(ctx, "user_notifications", "-", "+").Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения из Redis Stream")
		return err
	}

	for _, msg := range messages {
		if msg.Values["user_id"] == strconv.FormatInt(userID, 10) {
			_, err := c.client.XDel(ctx, "user_notifications", msg.ID).Result()
			if err != nil {
				log.Error().Err(err).Msgf("Ошибка удаления уведомления пользователя %d", userID)
				return err
			}
		}
	}
	return nil
}

func (c *Cache) GetScheduleUserNotifications() ([]redis.XStream, error) {
	ctx := context.Background()

	streams, err := c.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"user_notifications", "0"},
		Count:   10,
		Block:   0, // Блокируемся и ждём новую задачу
	}).Result()

	if err != nil {
		return nil, err
	}
	return streams, nil
}

func (c *Cache) GetUserNotificationTime(userID int64) (string, error) {
	ctx := context.Background()

	// Читаем уведомление из Redis Stream
	messages, err := c.client.XRange(ctx, "user_notifications", "-", "+").Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения уведомлений из Redis")
		return "", fmt.Errorf("Redis is not available")
	}

	// Ищем уведомление текущего пользователя
	for _, msg := range messages {
		if msg.Values["user_id"] == strconv.FormatInt(userID, 10) {
			return msg.Values["executeAt"].(string), nil
		}
	}

	return "", nil
}
