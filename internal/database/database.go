package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Database struct {
	db *pgxpool.Pool
}

func NewDatabase(db *pgxpool.Pool) *Database {
	return &Database{db: db}
}

// Инициализация PostgreSQL
func Init(url string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %v", err)
	}

	// Проверяем и создаём таблицы
	if err := ensureTables(dbpool); err != nil {
		return nil, fmt.Errorf("Ошибка инициализации таблиц: %v", err)
	}

	return dbpool, nil

}

func ensureTables(db *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
			tg_id INT NOT NULL,
            name TEXT NOT NULL,
			city TEXT NOT NULL,
			city_id TEXT NOT NULL,
			region TEXT,
			notification_time TEXT,
			state TEXT
        );
		CREATE INDEX IF NOT EXISTS idx_tg_id ON users(tg_id);
		CREATE TABLE IF NOT EXISTS weather (
    		city_id INT PRIMARY KEY,
    		forecast JSONB NOT NULL,
    		updated_at TIMESTAMP DEFAULT NOW());`,
	}

	for _, query := range queries {
		if _, err := db.Exec(context.Background(), query); err != nil {
			log.Info().Msgf("Ошибка выполнения запроса: %s, ошибка: %v", query, err)
			return err
		}

	}

	return nil
}
