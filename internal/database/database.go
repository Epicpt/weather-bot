package database

import (
	"context"
	"fmt"
	"weather-bot/internal/cache"

	"github.com/jackc/pgx/v5"
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
			federal_subject TEXT,
			notification_time TIME,
			last_weather_forecast TEXT,
			state TEXT
        );
		CREATE INDEX IF NOT EXISTS idx_tg_id ON users(tg_id);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(context.Background(), query); err != nil {
			log.Info().Msgf("Ошибка выполнения запроса: %s, ошибка: %v", query, err)
			return err
		}

	}

	return nil
}

// SaveUserToDB записывает или обновляет пользователя в БД
func (d *Database) SaveUserToDB(u *cache.User) error {
	var existingCity string
	err := d.db.QueryRow(context.Background(), "SELECT city FROM users WHERE tg_id = $1", u.TgID).Scan(&existingCity)

	if err == pgx.ErrNoRows {
		// Пользователь новый → вставляем
		_, err = d.db.Exec(context.Background(), "INSERT INTO users (tg_id, name, city, city_id, federal_subject, state) VALUES ($1, $2, $3, $4, $5, $6)", u.TgID, u.Name, u.City, u.CityID, u.FederalSubject, u.State)
		if err == nil {
			log.Info().Msgf("Новый пользователь в БД: tg_id=%d, name=%s, city=%s, city_id=%s,federal_subject=%v, state=%s", u.TgID, u.Name, u.City, u.CityID, u.FederalSubject, u.State)
		}
	} else if err == nil {
		// Пользователь уже есть → обновляем name и city, city_id (если изменились)
		_, err = d.db.Exec(context.Background(), "UPDATE users SET name = $1, city = $2, city_id = $3, federal_subject = $4, state = $5 WHERE tg_id = $6", u.Name, u.City, u.CityID, u.FederalSubject, u.State, u.TgID)
		if err == nil {
			log.Info().Msgf("Обновлён пользователь в БД: tg_id=%d, name=%s, city изменён с %s на %s, state=%s", u.TgID, u.Name, existingCity, u.City, u.State)
		}
	}

	if err != nil {
		log.Info().Msgf("Ошибка записи в БД: %v", err)
	}
	return err
}

func (d *Database) GetUserFromDB(userID int64) (*cache.User, error) {
	var user cache.User
	var federalSubject *string

	err := d.db.QueryRow(context.Background(), `
	SELECT tg_id, name, city, city_id, federal_subject, state
	FROM users
	WHERE tg_id = $1
`, userID).Scan(&user.TgID, &user.Name, &user.City, &user.CityID, &federalSubject, &user.State)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("ошибка получения пользователя из БД: %w", err)
	}

	// Если федеральный субъект не nil, присваиваем его пользователю
	if federalSubject != nil {
		user.FederalSubject = federalSubject
	}

	log.Info().Msgf("Пользователь получен из БД: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s", user.TgID, user.Name, user.City, user.CityID, user.State)

	return &user, nil
}
