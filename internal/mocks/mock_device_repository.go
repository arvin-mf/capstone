package mocks

import (
	"caps_influx/internal/repository"
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockDeviceRepository struct {
	mock.Mock
}

func (m *MockDeviceRepository) FindDevices(ctx context.Context) ([]repository.Device, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repository.Device), args.Error(1)
}

func (m *MockDeviceRepository) FindDeviceByClientID(ctx context.Context, cID string) (*repository.Device, error) {
	args := m.Called(ctx, cID)
	return args.Get(0).(*repository.Device), args.Error(1)
}

func (m *MockDeviceRepository) AddDevice(ctx context.Context, params repository.Device) (sql.Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockDeviceRepository) DeleteDevice(ctx context.Context, params repository.Device) (sql.Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockDeviceRepository) FindDevicesWithSubject(ctx context.Context) ([]repository.DeviceWithSubject, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repository.DeviceWithSubject), args.Error(1)
}

func (m *MockDeviceRepository) SetDeviceSubject(ctx context.Context, params repository.SetDeviceSubjectParam) (sql.Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockDeviceRepository) RemoveDeviceSubject(ctx context.Context, dID int64) (sql.Result, error) {
	args := m.Called(ctx, dID)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockDeviceRepository) UpdateDeviceStatus(ctx context.Context, params repository.Device) (sql.Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockDeviceRepository) SetDeviceStatusToRedis(ctx context.Context, key, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockDeviceRepository) GetDeviceStatusFromRedis(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(string), args.Error(1)
}
