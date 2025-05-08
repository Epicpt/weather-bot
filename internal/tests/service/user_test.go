package tests

import (
	"errors"
	"testing"
	"weather-bot/internal/app/services"
	"weather-bot/internal/mocks"
	"weather-bot/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestUserService_SaveUser(t *testing.T) {
	tests := []struct {
		name             string
		user             *models.User
		mockPrimaryErr   error
		mockSecondaryErr error
		expectedErr      bool
	}{
		{
			name:             "Both storages succeed",
			user:             &models.User{TgID: 1},
			mockPrimaryErr:   nil,
			mockSecondaryErr: nil,
			expectedErr:      false,
		},
		{
			name:             "Primary fails, secondary succeeds",
			user:             &models.User{TgID: 2},
			mockPrimaryErr:   errors.New("primary error"),
			mockSecondaryErr: nil,
			expectedErr:      false,
		},
		{
			name:             "Primary succeeds, secondary fails",
			user:             &models.User{TgID: 3},
			mockPrimaryErr:   nil,
			mockSecondaryErr: errors.New("secondary error"),
			expectedErr:      false,
		},
		{
			name:             "Both fail",
			user:             &models.User{TgID: 4},
			mockPrimaryErr:   errors.New("primary error"),
			mockSecondaryErr: errors.New("secondary error"),
			expectedErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryMock := mocks.NewCache(t)
			secondaryMock := mocks.NewDatabase(t)

			primaryMock.On("SaveUser", tt.user).Return(tt.mockPrimaryErr)
			secondaryMock.On("SaveUser", tt.user).Return(tt.mockSecondaryErr)

			service := services.InitUserService(primaryMock, secondaryMock)

			err := service.SaveUser(tt.user)

			if tt.expectedErr {
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

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name              string
		userID            int64
		mockPrimaryUser   *models.User
		mockPrimaryErr    error
		mockSecondaryUser *models.User
		mockSecondaryErr  error
		expectedUser      *models.User
		expectErr         bool
	}{
		{
			name:            "Primary succeeds",
			userID:          1,
			mockPrimaryUser: &models.User{TgID: 1, Name: "PrimaryUser"},
			expectedUser:    &models.User{TgID: 1, Name: "PrimaryUser"},
		},
		{
			name:              "Primary fails, secondary succeeds",
			userID:            2,
			mockPrimaryErr:    errors.New("primary error"),
			mockSecondaryUser: &models.User{TgID: 2, Name: "SecondaryUser"},
			expectedUser:      &models.User{TgID: 2, Name: "SecondaryUser"},
		},
		{
			name:             "Both fail",
			userID:           3,
			mockPrimaryErr:   errors.New("primary error"),
			mockSecondaryErr: errors.New("secondary error"),
			expectErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryMock := mocks.NewCache(t)
			secondaryMock := mocks.NewDatabase(t)

			primaryMock.On("GetUser", tt.userID).Return(tt.mockPrimaryUser, tt.mockPrimaryErr)
			if tt.mockPrimaryErr != nil {
				secondaryMock.On("GetUser", tt.userID).Return(tt.mockSecondaryUser, tt.mockSecondaryErr)
			}

			service := services.InitUserService(primaryMock, secondaryMock)

			user, err := service.GetUser(tt.userID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				var dualErr *services.DualStorageError
				assert.ErrorAs(t, err, &dualErr)
				assert.Equal(t, tt.mockPrimaryErr, dualErr.Primary)
				assert.Equal(t, tt.mockSecondaryErr, dualErr.Secondary)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			primaryMock.AssertExpectations(t)
			secondaryMock.AssertExpectations(t)

		})
	}
}
