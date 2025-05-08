package tests

import (
	"errors"
	"testing"
	"weather-bot/internal/app/services"
	"weather-bot/internal/mocks"
	"weather-bot/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestWeatherService_SaveWeather(t *testing.T) {
	tests := []struct {
		name             string
		id               int
		forecast         *models.ProcessedForecast
		mockPrimaryErr   error
		mockSecondaryErr error
		expectErr        bool
	}{
		{
			name:             "Both succeed",
			id:               1,
			forecast:         &models.ProcessedForecast{},
			mockPrimaryErr:   nil,
			mockSecondaryErr: nil,
			expectErr:        false,
		},
		{
			name:             "Primary fails, secondary succeeds",
			id:               2,
			forecast:         &models.ProcessedForecast{},
			mockPrimaryErr:   errors.New("primary error"),
			mockSecondaryErr: nil,
			expectErr:        false,
		},
		{
			name:             "Both fail",
			id:               3,
			forecast:         &models.ProcessedForecast{},
			mockPrimaryErr:   errors.New("primary error"),
			mockSecondaryErr: errors.New("secondary error"),
			expectErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryMock := mocks.NewCache(t)
			secondaryMock := mocks.NewDatabase(t)

			primaryMock.On("SaveWeather", tt.id, tt.forecast).Return(tt.mockPrimaryErr)
			secondaryMock.On("SaveWeather", tt.id, tt.forecast).Return(tt.mockSecondaryErr)

			service := services.InitWeatherService(primaryMock, secondaryMock)

			err := service.SaveWeather(tt.id, tt.forecast)

			if tt.expectErr {
				assert.Error(t, err)
				var dualErr *services.DualStorageError
				assert.ErrorAs(t, err, &dualErr)
				assert.Equal(t, tt.mockPrimaryErr, dualErr.Primary)
				assert.Equal(t, tt.mockSecondaryErr, dualErr.Secondary)
			} else {
				assert.NoError(t, err)
			}

			primaryMock.AssertExpectations(t)
			secondaryMock.AssertExpectations(t)
		})
	}
}

func TestWeatherService_GetWeather(t *testing.T) {
	tests := []struct {
		name                 string
		id                   int
		mockPrimaryWeather   *models.ProcessedForecast
		mockPrimaryErr       error
		mockSecondaryWeather *models.ProcessedForecast
		mockSecondaryErr     error
		expectedWeather      *models.ProcessedForecast
		expectErr            bool
	}{
		{
			name:                 "Primary succeeds",
			id:                   1,
			mockPrimaryWeather:   &models.ProcessedForecast{},
			mockPrimaryErr:       nil,
			mockSecondaryWeather: &models.ProcessedForecast{},
			mockSecondaryErr:     nil,
			expectedWeather:      &models.ProcessedForecast{},
			expectErr:            false,
		},
		{
			name:                 "Primary fails, secondary succeeds",
			id:                   2,
			mockPrimaryWeather:   &models.ProcessedForecast{},
			mockPrimaryErr:       errors.New("primary error"),
			mockSecondaryWeather: &models.ProcessedForecast{},
			mockSecondaryErr:     nil,
			expectedWeather:      &models.ProcessedForecast{},
			expectErr:            false,
		},
		{
			name:                 "Both fail",
			id:                   3,
			mockPrimaryWeather:   &models.ProcessedForecast{},
			mockPrimaryErr:       errors.New("primary error"),
			mockSecondaryWeather: &models.ProcessedForecast{},
			mockSecondaryErr:     errors.New("secondary error"),
			expectErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryMock := mocks.NewCache(t)
			secondaryMock := mocks.NewDatabase(t)

			primaryMock.On("GetWeather", tt.id).Return(tt.mockPrimaryWeather, tt.mockPrimaryErr)
			if tt.mockPrimaryErr != nil {
				secondaryMock.On("GetWeather", tt.id).Return(tt.mockSecondaryWeather, tt.mockSecondaryErr)
			}

			service := services.InitWeatherService(primaryMock, secondaryMock)

			weather, err := service.GetWeather(tt.id)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, weather)
				var dualErr *services.DualStorageError
				assert.ErrorAs(t, err, &dualErr)
				assert.Equal(t, tt.mockPrimaryErr, dualErr.Primary)
				assert.Equal(t, tt.mockSecondaryErr, dualErr.Secondary)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedWeather, weather)
			}

			primaryMock.AssertExpectations(t)
			if tt.mockPrimaryErr != nil {
				secondaryMock.AssertExpectations(t)
			}
		})
	}
}
