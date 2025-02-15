package cache

import "time"

type User struct {
	TgID                int64      `json:"tg_id"`
	Name                string     `json:"name"`
	City                string     `json:"city"`
	CityID              string     `json:"city_id"`
	FederalSubject      *string    `json:"federal_subject,omitempty"`
	NotificationTime    *time.Time `json:"notification_time,omitempty"`
	LastWeatherForecast *string    `json:"last_weather_forecast,omitempty"`
	State               string     `json:"state"`
}

func NewUser(tgID int64, name string) *User {
	return &User{
		TgID: tgID,
		Name: name,
	}
}
