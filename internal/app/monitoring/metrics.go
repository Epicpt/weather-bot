package monitoring

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Метрики для работы бота
	BotRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bot_requests_total",
		Help: "Общее количество запросов от пользователей",
	})

	BotErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bot_errors_total",
		Help: "Общее количество ошибок при обработке запросов от пользователей",
	})

	BotUniqueUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bot_active_users",
		Help: "Количество уникальных пользователей",
	})

	// Метрики погоды
	WeatherAPIRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_api_requests_total",
		Help: "Количество запросов к OpenWeather API",
	})

	WeatherAPIErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_api_errors_total",
		Help: "Количество ошибок при запросах к OpenWeather API",
	})

	WeatherCacheHitsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_cache_hits_total",
		Help: "Количество успешных запросов погоды из кэша",
	})

	WeatherCacheMissesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_cache_misses_total",
		Help: "Количество запросов погоды, которых не было в кэше",
	})

	WeatherUpdateTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_update_total",
		Help: "Количество обновлений погоды в хранилищах",
	})

	WeatherUpdateFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "weather_update_failed",
		Help: "Сколько раз не удалось обновить погоду в хранилищах",
	})

	// Метрики уведомлений
	NotificationsSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notifications_sent_total",
		Help: "Сколько уведомлений было отправлено пользователям",
	})

	NotificationsFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notifications_failed_total",
		Help: "Сколько уведомлений не удалось отправить",
	})

	// Метрики Redis
	RedisConnectionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redis_connection_errors_total",
		Help: "Количество ошибок подключения к Redis",
	})

	RedisErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redis_connection_errors",
		Help: "Общее количество ошибок Redis",
	})

	RedisQueueLength = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "redis_queue_length",
		Help: "Размер очереди уведомлений в Redis",
	})

	RedisCacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redis_cache_hits_total",
		Help: "Количество успешных запросов в Redis-кэш",
	})

	RedisCacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redis_cache_misses_total",
		Help: "Количество запросов, которых не было в кэше",
	})

	// Метрики БД
	DBErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_errors_total",
		Help: "Общее количество ошибок БД",
	})
)

var (
	activeUsers = make(map[int64]struct{})
	mu          sync.Mutex
)

func UpdateUniqueUsers(userID int64) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := activeUsers[userID]; !exists {
		activeUsers[userID] = struct{}{}
		BotUniqueUsers.Set(float64(len(activeUsers)))
	}
}
