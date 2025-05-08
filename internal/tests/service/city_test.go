package tests

import (
	"errors"
	"testing"
	"weather-bot/internal/app/handlers"
	"weather-bot/internal/app/services"
	"weather-bot/internal/mocks"
	"weather-bot/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestCityService_SaveCity(t *testing.T) {
	tests := []struct {
		name         string
		primaryErr   error
		secondaryErr error
		wantErr      bool
		wantDualErr  bool
	}{
		{
			name:         "Both succed",
			primaryErr:   nil,
			secondaryErr: nil,
			wantErr:      false,
		},
		{
			name:         "Primary fails, secondary succeeds",
			primaryErr:   errors.New("primary error"),
			secondaryErr: nil,
			wantErr:      false,
		},
		{
			name:         "Primary succeeds, secondary fails",
			primaryErr:   nil,
			secondaryErr: errors.New("secondary error"),
			wantErr:      false,
		},
		{
			name:         "Both fail",
			primaryErr:   errors.New("primary error"),
			secondaryErr: errors.New("secondary error"),
			wantErr:      true,
			wantDualErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryMock := mocks.NewCache(t)
			secondaryMock := mocks.NewDatabase(t)

			primaryMock.On("SaveCity", mock.Anything).Return(tt.primaryErr)
			secondaryMock.On("SaveCity", mock.Anything).Return(tt.secondaryErr)

			service := services.InitCityService(primaryMock, secondaryMock)
			err := service.SaveCity(models.City{ID: 1, Name: "Test"})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantDualErr {
				_, ok := err.(*services.DualStorageError)
				assert.True(t, ok, "expected DualStorageError, got %T", err)
			}
			primaryMock.AssertExpectations(t)
			secondaryMock.AssertExpectations(t)
		})
	}
}

func TestCityService_GetCities(t *testing.T) {
	tests := []struct {
		name            string
		primaryCities   []models.City
		primaryErr      error
		secondaryCities []models.City
		secondaryErr    error
		expectedCities  []models.City
		wantErr         bool
		wantDualErr     bool
	}{
		{
			name:           "Primary returns cities",
			primaryCities:  []models.City{{ID: 1, Name: "City1"}},
			primaryErr:     nil,
			expectedCities: []models.City{{ID: 1, Name: "City1"}},
			wantErr:        false,
		},
		{
			name:            "Primary fails, secondary returns cities",
			primaryErr:      errors.New("primary error"),
			secondaryCities: []models.City{{ID: 2, Name: "City2"}},
			secondaryErr:    nil,
			expectedCities:  []models.City{{ID: 2, Name: "City2"}},
			wantErr:         false,
		},
		{
			name:         "Both fail",
			primaryErr:   errors.New("primary error"),
			secondaryErr: errors.New("secondary error"),
			wantErr:      true,
			wantDualErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryMock := mocks.NewCache(t)
			secondaryMock := mocks.NewDatabase(t)

			primaryMock.On("GetCities", mock.Anything).Return(tt.primaryCities, tt.primaryErr)
			if tt.primaryErr != nil {
				secondaryMock.On("GetCities", mock.Anything).Return(tt.secondaryCities, tt.secondaryErr)
			}

			service := services.InitCityService(primaryMock, secondaryMock)

			cities, err := service.GetCities("TestCity")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCities, cities)
			}

			if tt.wantDualErr {
				_, ok := err.(*services.DualStorageError)
				assert.True(t, ok, "expected DualStorageError, got %T", err)
			}

			primaryMock.AssertExpectations(t)
			secondaryMock.AssertExpectations(t)
		})
	}
}

func TestCityService_LoadCities(t *testing.T) {
	tests := []struct {
		name          string
		inputCities   []models.City
		primaryErrs   map[int]error // индекс → ошибка Primary
		secondaryErrs map[int]error // индекс → ошибка Secondary
		expectedErr   bool
	}{
		{
			name:          "All save successfully",
			inputCities:   []models.City{{ID: 1, Name: "City1"}, {ID: 2, Name: "City2"}},
			primaryErrs:   map[int]error{},
			secondaryErrs: map[int]error{},
			expectedErr:   false,
		},
		{
			name:          "First city: only primary fails (secondary ok)",
			inputCities:   []models.City{{ID: 1, Name: "City1"}},
			primaryErrs:   map[int]error{0: errors.New("primary fail")},
			secondaryErrs: map[int]error{0: nil},
			expectedErr:   false,
		},
		{
			name:          "First city: both fail → return error",
			inputCities:   []models.City{{ID: 1, Name: "City1"}},
			primaryErrs:   map[int]error{0: errors.New("primary fail")},
			secondaryErrs: map[int]error{0: errors.New("secondary fail")},
			expectedErr:   true,
		},
		{
			name:          "Second city: both fail → return error",
			inputCities:   []models.City{{ID: 1, Name: "City1"}, {ID: 2, Name: "City2"}},
			primaryErrs:   map[int]error{1: errors.New("primary fail")},
			secondaryErrs: map[int]error{1: errors.New("secondary fail")},
			expectedErr:   true,
		},
		{
			name:          "Empty input list",
			inputCities:   []models.City{},
			primaryErrs:   map[int]error{},
			secondaryErrs: map[int]error{},
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryMock := mocks.NewCache(t)
			secondaryMock := mocks.NewDatabase(t)

			service := services.InitCityService(primaryMock, secondaryMock)

			for i, city := range tt.inputCities {
				pErr := tt.primaryErrs[i]
				sErr := tt.secondaryErrs[i]

				primaryMock.On("SaveCity", city).Return(pErr)
				secondaryMock.On("SaveCity", city).Return(sErr)
			}

			err := service.LoadCities(tt.inputCities)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			primaryMock.AssertExpectations(t)
			secondaryMock.AssertExpectations(t)
		})
	}
}

func TestCityService_GetCitiesNames(t *testing.T) {
	runGetTest := func(t *testing.T, getFunc func(s *services.CityService) ([]string, error), methodName string) {
		tests := []struct {
			name         string
			primaryRes   []string
			primaryErr   error
			secondaryRes []string
			secondaryErr error
			expectedRes  []string
			expectErr    bool
		}{
			{
				name:        "Primary success",
				primaryRes:  []string{"City1", "City2"},
				primaryErr:  nil,
				expectedRes: []string{"City1", "City2"},
				expectErr:   false,
			},
			{
				name:         "Primary fail, secondary success",
				primaryErr:   errors.New("primary fail"),
				secondaryRes: []string{"City3", "City4"},
				secondaryErr: nil,
				expectedRes:  []string{"City3", "City4"},
				expectErr:    false,
			},
			{
				name:         "Both fail",
				primaryErr:   errors.New("primary fail"),
				secondaryErr: errors.New("secondary fail"),
				expectErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				primaryMock := mocks.NewCache(t)
				secondaryMock := mocks.NewDatabase(t)

				service := services.InitCityService(primaryMock, secondaryMock)

				primaryMock.On(methodName).Return(tt.primaryRes, tt.primaryErr)
				if tt.primaryErr != nil {
					secondaryMock.On(methodName).Return(tt.secondaryRes, tt.secondaryErr)
				}

				result, err := getFunc(&service)

				if tt.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedRes, result)
				}

				primaryMock.AssertExpectations(t)
				secondaryMock.AssertExpectations(t)
			})
		}
	}

	runGetTest(t, func(s *services.CityService) ([]string, error) {
		return s.GetCitiesNames()
	}, "GetCitiesNames")
}

func TestCityService_GetCitiesIds(t *testing.T) {
	runGetTest := func(t *testing.T, getFunc func(s *services.CityService) ([]string, error), methodName string) {
		tests := []struct {
			name         string
			primaryRes   []string
			primaryErr   error
			secondaryRes []string
			secondaryErr error
			expectedRes  []string
			expectErr    bool
		}{
			{
				name:        "Primary success",
				primaryRes:  []string{"1", "2"},
				primaryErr:  nil,
				expectedRes: []string{"1", "2"},
				expectErr:   false,
			},
			{
				name:         "Primary fail, secondary success",
				primaryErr:   errors.New("primary fail"),
				secondaryRes: []string{"3", "4"},
				secondaryErr: nil,
				expectedRes:  []string{"3", "4"},
				expectErr:    false,
			},
			{
				name:         "Both fail",
				primaryErr:   errors.New("primary fail"),
				secondaryErr: errors.New("secondary fail"),
				expectErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				primaryMock := mocks.NewCache(t)
				secondaryMock := mocks.NewDatabase(t)

				service := services.InitCityService(primaryMock, secondaryMock)

				primaryMock.On(methodName).Return(tt.primaryRes, tt.primaryErr)
				if tt.primaryErr != nil {
					secondaryMock.On(methodName).Return(tt.secondaryRes, tt.secondaryErr)
				}

				result, err := getFunc(&service)

				if tt.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedRes, result)
				}

				primaryMock.AssertExpectations(t)
				secondaryMock.AssertExpectations(t)
			})
		}
	}

	runGetTest(t, func(s *services.CityService) ([]string, error) {
		return s.GetCitiesIds()
	}, "GetCitiesIds")
}
