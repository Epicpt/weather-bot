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

func TestRemoveUserNotification_Success(t *testing.T) {
	mockStorage := new(MockNotificationStorage)
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
	mockStorage := new(MockNotificationStorage)
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
	mockStorage := new(MockNotificationStorage)
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
	mockStorage := new(MockNotificationStorage)
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
	mockStorage := new(MockNotificationStorage)
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
	mockStorage := new(MockNotificationStorage)
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
	mockStorage := new(MockNotificationStorage)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	mockStorage.On("RemoveWeatherUpdate").Return(nil)

	err := service.RemoveWeatherUpdate()

	assert.NoError(t, err)

	mockStorage.AssertCalled(t, "RemoveWeatherUpdate")
}

func TestRemoveWeatherUpdate_Error(t *testing.T) {
	mockStorage := new(MockNotificationStorage)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	mockStorage.On("RemoveWeatherUpdate").Return(fmt.Errorf("storage error"))

	err := service.RemoveWeatherUpdate()

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "RemoveWeatherUpdate")
}

func TestGetScheduleWeatherUpdate_Success(t *testing.T) {
	mockStorage := new(MockNotificationStorage)
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
	mockStorage := new(MockNotificationStorage)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	mockStorage.On("GetScheduleWeatherUpdate").Return(nil, fmt.Errorf("storage error"))

	_, err := service.GetScheduleWeatherUpdate()

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "GetScheduleWeatherUpdate")
}

func TestGetScheduleUserNotifications_Success(t *testing.T) {
	mockStorage := new(MockNotificationStorage)
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
	mockStorage := new(MockNotificationStorage)
	service := &services.NotificationService{
		Primary: mockStorage,
	}

	mockStorage.On("GetScheduleUserNotifications").Return(nil, fmt.Errorf("storage error"))

	_, err := service.GetScheduleUserNotifications()

	assert.Error(t, err)

	mockStorage.AssertCalled(t, "GetScheduleUserNotifications")
}
