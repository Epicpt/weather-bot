package app

import (
	"os"
	"path/filepath"
	"weather-bot/internal/app/handlers"
	"weather-bot/internal/app/jobs"
	"weather-bot/internal/app/loader"
	"weather-bot/internal/app/reply"
	"weather-bot/internal/app/services"
	"weather-bot/internal/cache"
	"weather-bot/internal/config"
	"weather-bot/internal/database"
	"weather-bot/pkg/logger"
	"weather-bot/pkg/telegram"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Application interface {
	Bootstrap()
	Run()
	Shutdown()
}

type App struct {
	Bot   *tgbotapi.BotAPI
	DB    *database.Database
	Cache *cache.Cache
	Log   zerolog.Logger
}

func New(cfg *config.Config) *App {
	// Инициализация логгера
	log := logger.New()
	log.Info().Msg("Logger initialized")

	// Инициализация Postgres
	pool, err := database.Init(cfg.PostgresURL)
	if err != nil {
		log.Fatal().Err(err).Msgf("Ошибка подключения к БД: %v", err)
	}
	log.Info().Msg("Connected to PostgreSQL")

	// Инициализация Redis
	client, err := cache.Init(cfg.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка подключения к Redis")
	}

	log.Info().Msg("Connected to Redis")

	db := database.NewDatabase(pool)
	redis := cache.NewCache(client)

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating bot")
	}

	log.Info().Msgf("Бот %s запущен", bot.Self.UserName)

	return &App{
		Bot:   bot,
		DB:    db,
		Cache: redis,
		Log:   log,
	}
}

func (a *App) Bootstrap() {
	services.Init(a.Cache, a.DB)

	reply.Init(telegram.New(a.Bot))

	// Загрузка городов
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка получения текущего каталога")
	}
	filePath := filepath.Join(basePath, "internal", "app", "loader", "enriched_cities.json")
	if err := loader.LoadCities(filePath, services.InitCityService(a.Cache, a.DB)); err != nil {
		log.Fatal().Err(err).Msg("Error loading cities to storage")
	}

	log.Info().Msg("Cities loaded to Redis and Database")

	jobs.Init()
}

func (a *App) Run() {
	a.Log.Info().Msg("Bot started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := a.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Пропускаем неполные сообщения
			continue
		}

		handlers.Update(update)
	}
}

func (a *App) Shutdown() {
	log.Info().Msg("Отключение БД и Redis...")
	a.DB.Close()
	a.Cache.Close()
}
