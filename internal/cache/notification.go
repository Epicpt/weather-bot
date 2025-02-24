package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (c *Cache) SetNotificationTime(userID int, timeStr string) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("notify_time:user:%d", userID)

	err := c.c.Set(ctx, cacheKey, timeStr, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("ошибка сохранения времени уведомлений в Redis: %w", err)
	}

	return nil
}

func (c *Cache) GetNotificationTime(userID int) (string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("notify_time:user:%d", userID)

	timeStr, err := c.c.Get(ctx, cacheKey).Result()
	if err == nil {
		return timeStr, nil
	}

	return timeStr, nil
}

func (c *Cache) SetUpdateWeather(executeAt int64) error {
	ctx := context.Background()
	// ID задачи
	jobID := fmt.Sprintf("weather_update:%d", executeAt)

	// Добавляем задачу в Redis Stream
	_, err := c.c.XAdd(ctx, &redis.XAddArgs{
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

func (c *Cache) DeleteUpdateWeather() error {
	ctx := context.Background()

	messages, err := c.c.XRange(ctx, "weather_updates", "-", "+").Result()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения из Redis Stream")
		return err
	}

	for _, msg := range messages {
		_, err := c.c.XDel(ctx, "weather_updates", msg.ID).Result()
		if err != nil {
			log.Error().Err(err).Msgf("Ошибка удаления задачи %s", msg.ID)
			return err
		}
	}
	return nil
}

func (c *Cache) GetUpdateWeather() ([]redis.XStream, error) {
	ctx := context.Background()

	streams, err := c.c.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"weather_updates", "0"},
		Count:   1,
		Block:   0, // Блокируемся и ждём новую задачу
	}).Result()

	if err != nil {
		return nil, err
	}
	return streams, nil
}
