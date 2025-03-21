package tests

import (
	"testing"
	"weather-bot/internal/app/handlers"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCity(t *testing.T) {
	tests := []struct {
		city     string
		expected bool
	}{
		{"Москва", true},
		{"Санкт-Петербург", true},
		{"Новосибирск", true},
		{"Мурманск", true},
		{"Томск", true},
		{"Омск", true},
		{"Красноярск", true},
		{"Ростов-на-Дону", true},
		{"Самара", true},
		{"Казань", true},
		{"Екатеринбург", true},
		{"Нижний Новгород", true},
		{"Владимир", true},
		{"Уфа", true},
		{"Пермь", true},
		{"Сочи", true},
		{"Челябинск", true},
		{"City", false},
		{"123", false},
		{"Москва!", false},
		{"", false},
		{"-Казань-", true},
		{"Москва 123", false},
		{"/city", false},
	}

	for _, tt := range tests {
		t.Run(tt.city, func(t *testing.T) {
			assert.Equal(t, tt.expected, handlers.IsValidCity(tt.city), "Город: %s", tt.city)
		})
	}

}
