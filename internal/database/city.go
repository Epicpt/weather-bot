package database

import (
	"context"
)

func (d *Database) GetCitiesIds() ([]string, error) {
	ctx := context.Background()
	var cityIDs []string

	rows, err := d.db.Query(ctx, "SELECT DISTINCT city_id FROM users")
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
