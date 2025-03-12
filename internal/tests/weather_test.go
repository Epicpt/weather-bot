package tests

import (
	"fmt"
	"testing"
	"weather-bot/internal/app/services"
	"weather-bot/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherStorage struct {
	mock.Mock
}

func (m *MockWeatherStorage) SaveWeather(id int, forecast *models.ProcessedForecast) error {
	args := m.Called(id, forecast)
	return args.Error(0)
}

func (m *MockWeatherStorage) GetWeather(id int) (*models.ProcessedForecast, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProcessedForecast), args.Error(1)
}

func TestSaveWeather_PrimaryFailsSecondaryFails(t *testing.T) {
	mockPrimary := new(MockWeatherStorage)
	mockSecondary := new(MockWeatherStorage)
	service := &services.WeatherService{
		Primary:   mockPrimary,
		Secondary: mockSecondary,
	}

	weather := &models.ProcessedForecast{}
	// Мокируем ошибку для первичного хранилища
	mockPrimary.On("SaveWeather", 1, weather).Return(fmt.Errorf("primary error"))
	// Мокируем ошибку для вторичного хранилища
	mockSecondary.On("SaveWeather", 1, weather).Return(fmt.Errorf("secondary error"))

	// Вызываем метод SaveWeather
	err := service.SaveWeather(1, weather)

	// Проверяем, что ошибка не nil
	assert.Error(t, err)

	// Приводим ошибку к типу DualStorageError и проверяем её поля
	var dualErr = &services.DualStorageError{
		Primary:   fmt.Errorf("primary error"),
		Secondary: fmt.Errorf("secondary error"),
	}
	if assert.IsType(t, dualErr, err) {
		assert.Equal(t, "Primary: primary error, Secondary: secondary error", dualErr.Error()) // проверяем ошибку в первичном хранилище
		assert.Equal(t, "Primary: primary error, Secondary: secondary error", dualErr.Error()) // проверяем ошибку во вторичном хранилище
	}

	// Проверяем, что методы действительно были вызваны
	mockPrimary.AssertCalled(t, "SaveWeather", 1, weather)
	mockSecondary.AssertCalled(t, "SaveWeather", 1, weather)
}

func TestSaveWeather_PrimaryFailsSecondarySucceeds(t *testing.T) {
	mockPrimary := new(MockWeatherStorage)
	mockSecondary := new(MockWeatherStorage)
	service := &services.WeatherService{
		Primary:   mockPrimary,
		Secondary: mockSecondary,
	}

	weather := &models.ProcessedForecast{}
	// Мокируем ошибку для первичного хранилища
	mockPrimary.On("SaveWeather", 1, weather).Return(fmt.Errorf("primary error"))
	// Мокируем успешное сохранение во вторичном хранилище
	mockSecondary.On("SaveWeather", 1, weather).Return(nil)

	// Вызываем метод SaveWeather
	err := service.SaveWeather(1, weather)

	// Проверяем, что ошибка не nil
	assert.NoError(t, err)

	// Проверяем, что методы действительно были вызваны
	mockPrimary.AssertCalled(t, "SaveWeather", 1, weather)
	mockSecondary.AssertCalled(t, "SaveWeather", 1, weather)
}

func TestSaveWeather_PrimarySucceedsSecondaryFails(t *testing.T) {
	mockPrimary := new(MockWeatherStorage)
	mockSecondary := new(MockWeatherStorage)
	service := &services.WeatherService{
		Primary:   mockPrimary,
		Secondary: mockSecondary,
	}

	weather := &models.ProcessedForecast{}
	// Мокируем успешное сохранение в первичном хранилище
	mockPrimary.On("SaveWeather", 1, weather).Return(nil)
	// Мокируем ошибку для вторичного хранилища
	mockSecondary.On("SaveWeather", 1, weather).Return(fmt.Errorf("secondary error"))

	// Вызываем метод SaveWeather
	err := service.SaveWeather(1, weather)

	// Проверяем, что ошибка не произошла (ошибка только во вторичном хранилище)
	assert.NoError(t, err)

	// Проверяем, что методы действительно были вызваны
	mockPrimary.AssertCalled(t, "SaveWeather", 1, weather)
	mockSecondary.AssertCalled(t, "SaveWeather", 1, weather)
}
