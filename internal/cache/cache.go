package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
	"weather-bot/internal/app/monitoring"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	client  *redis.Client
	Healthy bool
	mu      sync.RWMutex
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

// Инициализация Redis
func Init(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0, // use default DB
	})

	if client == nil {
		return nil, fmt.Errorf("Redis клиент не инициализирован")
	}

	return client, nil

}

func (cache *Cache) HealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := cache.client.Ping(ctx).Result()
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if err != nil {
		monitoring.RedisConnectionErrors.Inc()
		log.Warn().Msg("Redis недоступен")
		cache.Healthy = false
	} else {
		if !cache.Healthy {
			log.Info().Msg("Redis доступен")

		}
		cache.Healthy = true
	}
}

func (cache *Cache) IsHealthy() bool {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.Healthy
}

func (cache *Cache) Close() {
	cache.client.Close()
}
