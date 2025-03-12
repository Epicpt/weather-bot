package tests

import (
	"fmt"
	"testing"
	"weather-bot/internal/app/services"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotificationStorage struct {
	mock.Mock
}

func (m *MockNotificationStorage) GetUserNotificationTime(userID int64) (string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *MockNotificationStorage) RemoveUserNotification(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockNotificationStorage) ScheduleUserNotification(userID int64, executeAt int64) error {
	args := m.Called(userID, executeAt)
	return args.Error(0)
}

func (m *MockNotificationStorage) ScheduleWeatherUpdate(executeAt int64) error {
	args := m.Called(executeAt)
	return args.Error(0)
}

func (m *MockNotificationStorage) RemoveWeatherUpdate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNotificationStorage) GetScheduleWeatherUpdate() ([]redis.XStream, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]redis.XStream), args.Error(1)
}
func (m *MockNotificationStorage) GetScheduleUserNotifications() ([]redis.XStream, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]redis.XStream), args.Error(1)
}

func TestGetUserNotificationTime_Success(t *testing.T) {
	mockStorage := new(MockNotificationStorage)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	userID := int64(1)
	expectedTime := "08:00"

	mockStorage.On("GetUserNotificationTime", userID).Return(expectedTime, nil)

	notifTime, err := service.GetUserNotificationTime(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTime, notifTime)

	mockStorage.AssertCalled(t, "GetUserNotificationTime", userID)
}

func TestGetUserNotificationTime_Error(t *testing.T) {
	mockStorage := new(MockNotificationStorage)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	userID := int64(1)

	mockStorage.On("GetUserNotificationTime", userID).Return("", fmt.Errorf("storage error"))

	notifTime, err := service.GetUserNotificationTime(userID)

	assert.Error(t, err)
	assert.Empty(t, notifTime)

	mockStorage.AssertCalled(t, "GetUserNotificationTime", userID)
}
