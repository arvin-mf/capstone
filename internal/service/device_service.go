package service

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/repository"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
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

	deviceCount := len(devices)
	errChan := make(chan error, deviceCount)

	var wg sync.WaitGroup
	wg.Add(deviceCount)

	for _, device := range devices {
		deviceCopy := device

		go func(d *repository.DeviceWithSubject) {
			defer wg.Done()

			if d.DeviceStatus {
				key := "device-" + strconv.Itoa(int(d.DeviceID))
				_, err := s.deviceRepo.GetDeviceStatusFromRedis(ctx, key)
				if err != nil {
					if !errors.Is(err, redis.Nil) {
						errChan <- err
						return
					} else {
						_, er := s.deviceRepo.UpdateDeviceStatus(ctx, repository.Device{
							ID:           d.DeviceID,
							DeviceStatus: false,
						})
						if er != nil {
							errChan <- err
							return
						}
					}
				}
			}
		}(&deviceCopy)
	}
	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("Failed to update device status: %v\n", e)
		}
		return nil, errors.New("failed to update device status")
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
