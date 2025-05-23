package service

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/repository"
	"context"
)

type DeviceService interface {
	GetDevices(ctx context.Context) ([]dto.DeviceResponse, error)
	AddDevice(ctx context.Context, arg dto.DeviceCreateReq) error
	DeleteDevice(ctx context.Context, deviceID int64) error
	GetDevicesWithSubject(ctx context.Context) ([]dto.DeviceWithSubjectResponse, error)
	SetDeviceSubject(ctx context.Context, deviceID int64, arg dto.SetDeviceSubjectReq) error
}

type deviceService struct {
	deviceRepo repository.DeviceRepository
}

func NewDeviceService(dr repository.DeviceRepository) DeviceService {
	return &deviceService{
		deviceRepo: dr,
	}
}

func (s *deviceService) GetDevices(ctx context.Context) ([]dto.DeviceResponse, error) {
	devices, err := s.deviceRepo.FindDevices(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToDeviceResponses(&devices), nil
}

func (s *deviceService) AddDevice(ctx context.Context, arg dto.DeviceCreateReq) error {
	_, err := s.deviceRepo.AddDevice(ctx, repository.Device{
		ClientID: arg.ClientID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *deviceService) DeleteDevice(ctx context.Context, deviceID int64) error {
	_, err := s.deviceRepo.DeleteDevice(ctx, repository.Device{
		ID: deviceID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *deviceService) GetDevicesWithSubject(ctx context.Context) ([]dto.DeviceWithSubjectResponse, error) {
	devices, err := s.deviceRepo.FindDevicesWithSubject(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToDeviceWithSubjectResponses(&devices), nil
}

func (s *deviceService) SetDeviceSubject(ctx context.Context, deviceID int64, arg dto.SetDeviceSubjectReq) error {
	_, err := s.deviceRepo.RemoveDeviceSubject(ctx, deviceID)
	if err != nil {
		return err
	}

	_, err = s.deviceRepo.SetDeviceSubject(ctx, repository.SetDeviceSubjectParam{
		SubjectID: arg.SubjectID,
		DeviceID:  deviceID,
	})
	if err != nil {
		return err
	}

	return nil
}
