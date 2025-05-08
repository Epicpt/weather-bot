package services

import (
	"weather-bot/internal/app/monitoring"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/rs/zerolog/log"
)

type UserService struct {
	Primary   storage.UserStorage
	Secondary storage.UserStorage
}

func (s *UserService) SaveUser(user *models.User) error {
	var errP, errS error
	errP = s.Primary.SaveUser(user)
	if errP != nil {
		monitoring.RedisErrorsTotal.Inc()
		log.Warn().Err(errP).Msg("Ошибка записи юзера в Primary хранилище")
	}

	if errS = s.Secondary.SaveUser(user); errS != nil {
		monitoring.DBErrorsTotal.Inc()
		log.Warn().Err(errS).Msg("Ошибка записи юзера в Secondary хранилище")

	}

	if errP != nil && errS != nil {
		return &DualStorageError{Primary: errP, Secondary: errS}
	}

	return nil
}

func (s *UserService) GetUser(id int64) (*models.User, error) {
	user, errP := s.Primary.GetUser(id)
	if errP == nil {
		monitoring.RedisCacheHits.Inc()
		return user, nil
	}
	monitoring.RedisCacheMisses.Inc()
	monitoring.RedisErrorsTotal.Inc()
	log.Warn().Err(errP).Msg("Ошибка чтения юзера из Primary хранилища")

	user, errS := s.Secondary.GetUser(id)
	if errS == nil {
		return user, nil
	}
	monitoring.DBErrorsTotal.Inc()
	log.Warn().Err(errS).Msg("Ошибка чтения юзера из Secondary хранилища")

	return nil, &DualStorageError{Primary: errP, Secondary: errS}
}
