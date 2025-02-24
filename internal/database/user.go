package database

import (
	"context"
	"fmt"
	"weather-bot/internal/cache"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// SaveUserToDB записывает или обновляет пользователя в БД
func (d *Database) SaveUserToDB(u *cache.User) error {
	var existingCity string
	err := d.db.QueryRow(context.Background(), "SELECT city FROM users WHERE tg_id = $1", u.TgID).Scan(&existingCity)

	if err == pgx.ErrNoRows {
		// Пользователь новый → вставляем
		_, err = d.db.Exec(context.Background(), "INSERT INTO users (tg_id, name, city, city_id, region, state) VALUES ($1, $2, $3, $4, $5, $6)", u.TgID, u.Name, u.City, u.CityID, u.Region, u.State)
		if err == nil {
			log.Info().Msgf("Новый пользователь в БД: tg_id=%d, name=%s, city=%s, city_id=%s,region=%v, state=%s", u.TgID, u.Name, u.City, u.CityID, u.Region, u.State)
		}
	} else if err == nil {
		// Пользователь уже есть → обновляем name и city, city_id (если изменились)
		_, err = d.db.Exec(context.Background(), "UPDATE users SET name = $1, city = $2, city_id = $3, region = $4, state = $5 WHERE tg_id = $6", u.Name, u.City, u.CityID, u.Region, u.State, u.TgID)
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
	SELECT tg_id, name, city, city_id, region, state
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
		user.Region = federalSubject
	}

	log.Info().Msgf("Пользователь получен из БД: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s", user.TgID, user.Name, user.City, user.CityID, user.State)

	return &user, nil
}
