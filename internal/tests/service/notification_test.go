package tests

import (
	"fmt"
	"testing"
	"weather-bot/internal/app/services"
	"weather-bot/internal/mocks"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGetUserNotificationTime_Success(t *testing.T) {
	mockStorage := mocks.NewCache(t)
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
	mockStorage := mocks.NewCache(t)
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

func TestRemoveUserNotification_Success(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	userID := int64(1)

	mockStorage.On("RemoveUserNotification", userID).Return(nil)

	err := service.RemoveUserNotification(userID)

	assert.NoError(t, err)

	mockStorage.AssertCalled(t, "RemoveUserNotification", userID)
}

func TestRemoveUserNotification_Error(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	userID := int64(1)

	mockStorage.On("RemoveUserNotification", userID).Return(fmt.Errorf("storage error"))

	err := service.RemoveUserNotification(userID)

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "RemoveUserNotification", userID)
}

func TestScheduleUserNotification_Success(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	userID := int64(1)
	executeAt := int64(111)

	mockStorage.On("ScheduleUserNotification", userID, executeAt).Return(nil)

	err := service.ScheduleUserNotification(userID, executeAt)

	assert.NoError(t, err)

	mockStorage.AssertCalled(t, "ScheduleUserNotification", userID, executeAt)
}

func TestScheduleUserNotification_Error(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	userID := int64(1)
	executeAt := int64(111)

	mockStorage.On("ScheduleUserNotification", userID, executeAt).Return(fmt.Errorf("storage error"))

	err := service.ScheduleUserNotification(userID, executeAt)

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "ScheduleUserNotification", userID, executeAt)
}

func TestScheduleWeatherUpdate_Success(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	executeAt := int64(111)

	mockStorage.On("ScheduleWeatherUpdate", executeAt).Return(nil)

	err := service.ScheduleWeatherUpdate(executeAt)

	assert.NoError(t, err)

	mockStorage.AssertCalled(t, "ScheduleWeatherUpdate", executeAt)
}

func TestScheduleWeatherUpdate_Error(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	executeAt := int64(111)

	mockStorage.On("ScheduleWeatherUpdate", executeAt).Return(fmt.Errorf("storage error"))

	err := service.ScheduleWeatherUpdate(executeAt)

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "ScheduleWeatherUpdate", executeAt)
}

func TestRemoveWeatherUpdate_Success(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	mockStorage.On("RemoveWeatherUpdate").Return(nil)

	err := service.RemoveWeatherUpdate()

	assert.NoError(t, err)

	mockStorage.AssertCalled(t, "RemoveWeatherUpdate")
}

func TestRemoveWeatherUpdate_Error(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	mockStorage.On("RemoveWeatherUpdate").Return(fmt.Errorf("storage error"))

	err := service.RemoveWeatherUpdate()

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "RemoveWeatherUpdate")
}

func TestGetScheduleWeatherUpdate_Success(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	expectedSchedule := []redis.XStream{}

	mockStorage.On("GetScheduleWeatherUpdate").Return(expectedSchedule, nil)

	schedule, err := service.GetScheduleWeatherUpdate()

	assert.NoError(t, err)
	assert.Equal(t, expectedSchedule, schedule)

	mockStorage.AssertCalled(t, "GetScheduleWeatherUpdate")
}

func TestGetScheduleWeatherUpdate_Error(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	mockStorage.On("GetScheduleWeatherUpdate").Return(nil, fmt.Errorf("storage error"))

	_, err := service.GetScheduleWeatherUpdate()

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "GetScheduleWeatherUpdate")
}

func TestGetScheduleUserNotifications_Success(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	expectedSchedule := []redis.XStream{}

	mockStorage.On("GetScheduleUserNotifications").Return(expectedSchedule, nil)

	schedule, err := service.GetScheduleUserNotifications()

	assert.NoError(t, err)
	assert.Equal(t, expectedSchedule, schedule)

	mockStorage.AssertCalled(t, "GetScheduleUserNotifications")
}

func TestGetScheduleUserNotifications_Error(t *testing.T) {
	mockStorage := mocks.NewCache(t)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	mockStorage.On("GetScheduleUserNotifications").Return(nil, fmt.Errorf("storage error"))

	_, err := service.GetScheduleUserNotifications()

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "GetScheduleUserNotifications")
}
