package service

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/mocks"
	"caps_influx/internal/repository"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDevices(t *testing.T) {
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	service := NewDeviceService(mockDeviceRepo)

	devices := []repository.Device{
		{ID: 1, ClientID: "aa:bb:cc:dd:ee:ff", CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, ClientID: "bb:cc:dd:ee:ff:gg", CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 3, ClientID: "cc:dd:ee:ff:gg:hh", CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
	}
	mockDeviceRepo.On("FindDevices", mock.Anything).Return(devices, nil)

	result, err := service.GetDevices(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, dto.ToDeviceResponses(&devices), result)
	mockDeviceRepo.AssertExpectations(t)
}

func TestGetDevicesEmpty(t *testing.T) {
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	service := NewDeviceService(mockDeviceRepo)

	emptyDevices := []repository.Device{}
	mockDeviceRepo.On("FindDevices", mock.Anything).Return(emptyDevices, nil)

	result, err := service.GetDevices(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, dto.ToDeviceResponses(&emptyDevices), result)
	mockDeviceRepo.AssertExpectations(t)
}

func TestAddDevice(t *testing.T) {
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	service := NewDeviceService(mockDeviceRepo)

	device := dto.DeviceCreateReq{ClientID: "aa:bb:cc:dd:ee:ff"}
	mockResult := mocks.MockResult{}
	mockDeviceRepo.On("AddDevice", mock.Anything, mock.AnythingOfType("repository.Device")).Return(mockResult, nil)

	err := service.AddDevice(context.Background(), device)
	assert.NoError(t, err)
	mockDeviceRepo.AssertExpectations(t)
}

func TestDeleteDevice(t *testing.T) {
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	service := NewDeviceService(mockDeviceRepo)

	deviceID := int64(1)
	mockResult := mocks.MockResult{}
	mockDeviceRepo.On("DeleteDevice", mock.Anything, mock.AnythingOfType("repository.Device")).Return(mockResult, nil)

	err := service.DeleteDevice(context.Background(), deviceID)
	assert.NoError(t, err)
	mockDeviceRepo.AssertExpectations(t)
}

func TestGetDevicesWithSubject(t *testing.T) {
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	service := NewDeviceService(mockDeviceRepo)

	devicesWithSubject := []repository.DeviceWithSubject{
		{DeviceID: 1, SubjectID: sql.NullInt64{Int64: 10, Valid: true}, Name: sql.NullString{String: "Subject A", Valid: true}, IsFatigued: sql.NullBool{Bool: false, Valid: true}, CreatedAt: sql.NullTime{Time: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Valid: true}},
		{DeviceID: 2, SubjectID: sql.NullInt64{Int64: 11, Valid: true}, Name: sql.NullString{String: "Subject B", Valid: true}, IsFatigued: sql.NullBool{Bool: false, Valid: true}, CreatedAt: sql.NullTime{Time: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Valid: true}},
	}
	mockDeviceRepo.On("FindDevicesWithSubject", mock.Anything).Return(devicesWithSubject, nil)

	result, err := service.GetDevicesWithSubject(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, dto.ToDeviceWithSubjectResponses(&devicesWithSubject), result)
	mockDeviceRepo.AssertExpectations(t)
}

func TestGetDevicesWithSubjectEmpty(t *testing.T) {
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	service := NewDeviceService(mockDeviceRepo)

	emptyDevicesWithSubject := []repository.DeviceWithSubject{}
	mockDeviceRepo.On("FindDevicesWithSubject", mock.Anything).Return(emptyDevicesWithSubject, nil)

	result, err := service.GetDevicesWithSubject(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, dto.ToDeviceWithSubjectResponses(&emptyDevicesWithSubject), result)
	mockDeviceRepo.AssertExpectations(t)
}

func TestSetDeviceSubject(t *testing.T) {
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	service := NewDeviceService(mockDeviceRepo)

	deviceID := int64(1)
	arg := dto.SetDeviceSubjectReq{SubjectID: 10}
	mockResult := mocks.MockResult{}
	mockDeviceRepo.On("RemoveDeviceSubject", mock.Anything, deviceID).Return(mockResult, nil)
	mockDeviceRepo.On("SetDeviceSubject", mock.Anything, repository.SetDeviceSubjectParam{
		SubjectID: arg.SubjectID,
		DeviceID:  deviceID,
	}).Return(mockResult, nil)

	err := service.SetDeviceSubject(context.Background(), deviceID, arg)
	assert.NoError(t, err)
	mockDeviceRepo.AssertExpectations(t)
}
