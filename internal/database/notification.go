package database

import (
	"context"
	"fmt"
)

func (d *Database) SetNotificationTime(userID int, timeStr string) error {
	ctx := context.Background()
	_, err := d.db.Exec(ctx, "UPDATE users SET notification_time = $1 WHERE tg_id = $2", timeStr, userID)
	if err != nil {
		return fmt.Errorf("ошибка сохранения времени уведомлений в БД: %w", err)
	}

	return nil
}

func (d *Database) GetNotificationTime(userID int) (string, error) {
	ctx := context.Background()
	var timeStr string

	err := d.db.QueryRow(ctx, "SELECT notification_time FROM users WHERE tg_id = $1", userID).Scan(&timeStr)
	if err != nil {
		return "", fmt.Errorf("не удалось получить время уведомлений: %w", err)
	}

	return timeStr, nil
}
