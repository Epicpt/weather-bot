package database

import (
	"context"
	"weather-bot/internal/app/storage"
	"weather-bot/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var _ storage.CityStorage = (*Database)(nil) // Проверка интерфейса

func (db *Database) SaveCity(city models.City) error {

	_, err := db.pool.Exec(context.Background(), `
			INSERT INTO cities (id, name, federal_district, region, city_district, street)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE SET name = $2, federal_district = $3, region = $4, city_district = $5, street = $6`,
		city.ID, city.Name, city.FederalDistrict, city.Region, city.CityDistrict, city.Street,
	)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка записи города в БД")
		return err
	}

	return nil
}

// GetCities ищет города в PostgreSQL по имени
func (db *Database) GetCities(name string) ([]models.City, error) {
	ctx := context.Background()

	rows, err := db.pool.Query(ctx, `
		SELECT id, name, federal_district, region, city_district, street 
		FROM cities 
		WHERE name = $1`, name)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка запроса к БД")
		return nil, err
	}
	defer rows.Close()

	seen := make(map[string]bool)
	var result []models.City

	for rows.Next() {
		var city models.City
		err := rows.Scan(&city.ID, &city.Name, &city.FederalDistrict, &city.Region, &city.CityDistrict, &city.Street)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка чтения данных из БД")
			continue
		}

		// Исключаем дубли по `Region`
		if seen[city.Region] {
			continue
		}
		seen[city.Region] = true

		result = append(result, city)
	}

	if len(result) == 0 {
		return nil, pgx.ErrNoRows
	}

	return result, nil
}

func (d *Database) GetCitiesIds() ([]string, error) {
	ctx := context.Background()
	var cityIDs []string

	rows, err := d.pool.Query(ctx, "SELECT DISTINCT city_id FROM users")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cityID string
		if err := rows.Scan(&cityID); err == nil {
			cityIDs = append(cityIDs, cityID)
		}
	}

	return cityIDs, nil
}

func (d *Database) GetCitiesNames() ([]string, error) {
	ctx := context.Background()
	var citiesNames []string

	rows, err := d.pool.Query(ctx, "SELECT DISTINCT name FROM cities")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cityName string
		if err := rows.Scan(&cityName); err == nil {
			citiesNames = append(citiesNames, cityName)
		}
	}

	return citiesNames, nil
}
