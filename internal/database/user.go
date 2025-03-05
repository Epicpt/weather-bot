package database

import (
	"context"
	"fmt"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var _ storage.UserStorage = (*Database)(nil)

// SaveUser записывает или обновляет пользователя в БД
func (d *Database) SaveUser(u *models.User) error {
	_, err := d.pool.Exec(context.Background(), `
		INSERT INTO users (tg_id, chat_id, name, city, city_id, region, state, sticker) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		ON CONFLICT (tg_id) DO UPDATE SET chat_id = $2, name = $3, city = $4, city_id = $5, region = $6, state = $7, sticker = $8`,
		u.TgID, u.ChatID, u.Name, u.City, u.CityID, u.Region, u.State, u.Sticker)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка записи юзера в БД")
	}

	log.Info().Msgf("Пользователь сохранён в БД: tg_id=%d, chat_id=%d, name=%s, city=%s, city_id=%s, region=%v, state=%s, sticker=%v", u.TgID, u.ChatID, u.Name, u.City, u.CityID, u.Region, u.State, u.Sticker)

	return nil
}

func (d *Database) GetUser(userID int64) (*models.User, error) {
	var user models.User

	err := d.pool.QueryRow(context.Background(), `
	SELECT tg_id, name, city, city_id, region, state, sticker
	FROM users
	WHERE tg_id = $1
`, userID).Scan(&user.TgID, &user.Name, &user.City, &user.CityID, &user.Region, &user.State, &user.Sticker)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("ошибка получения пользователя из БД: %w", err)
	}

	log.Info().Msgf("Пользователь получен из БД: tg_id=%d, name=%s, city=%s, city_id=%s, state=%s, sticker=%v", user.TgID, user.Name, user.City, user.CityID, user.State, user.Sticker)

	return &user, nil
}
