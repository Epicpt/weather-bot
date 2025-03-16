package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

// Инициализация PostgreSQL
func Init(url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %v", err)
	}

	// Проверяем и создаём таблицы
	if err := ensureTables(pool); err != nil {
		return nil, fmt.Errorf("Ошибка инициализации таблиц: %v", err)
	}

	return pool, nil

}

func ensureTables(pool *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			tg_id INT PRIMARY KEY,
			chat_id INT NOT NULL,
            name TEXT NOT NULL,
			city TEXT NOT NULL,
			city_id TEXT NOT NULL,
			region TEXT,
			state TEXT,
			sticker BOOLEAN DEFAULT FALSE
        );
		CREATE INDEX IF NOT EXISTS idx_tg_id ON users(tg_id);
		CREATE TABLE IF NOT EXISTS weather (
    		city_id INT PRIMARY KEY,
    		forecast JSONB NOT NULL,
    		updated_at TIMESTAMP DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS cities (
    		id SERIAL PRIMARY KEY,
    		name TEXT NOT NULL,
    		federal_district TEXT,
    		region TEXT,
    		city_district TEXT,
    		street TEXT
		);`,
	}

	for _, query := range queries {
		if _, err := pool.Exec(context.Background(), query); err != nil {
			log.Debug().Msgf("Ошибка выполнения запроса: %s, ошибка: %v", query, err)
			return err
		}

	}

	return nil
}

func (db *Database) Close() {
	db.pool.Close()
}
