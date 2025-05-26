package service

import (
	"caps_influx/internal/mocks"
	"caps_influx/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestWritePeriodicDeviceExists(t *testing.T) {
	mockInfluxRepo := new(mocks.MockInfluxRepository)
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	mockSubjectRepo := new(mocks.MockSubjectRepository)
	service := NewInfluxService(mockInfluxRepo, mockSubjectRepo, mockDeviceRepo)

	mockMessage := &mocks.MockMessage{
		MockPayload: []byte(`{"bpm":62,"suhu_objek":36.5,"suhu_lingkungan":22.0,"status":1}`),
		MockTopic:   "esp32/aa:bb:cc:dd:ee:ff/periodic",
	}

	mockDevice := &repository.Device{ID: 7, ClientID: "aa:bb:cc:dd:ee:ff"}
	mockSubject := &repository.Subject{ID: 10, Name: "My Subject", IsFatigued: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockResult := &mocks.MockResult{}

	mockDeviceRepo.On("FindDeviceByClientID", mock.Anything, "aa:bb:cc:dd:ee:ff").Return(mockDevice, nil)
	mockSubjectRepo.On("FindSubjectByDeviceID", mock.Anything, int64(7)).Return(mockSubject, nil)
	mockSubjectRepo.On("UpdateSubjectFatiguedStatus", mock.Anything, repository.Subject{
		ID:         10,
		IsFatigued: true,
	}).Return(mockResult, nil)

	expectedParam := repository.InfluxPeriodicPointParam{
		DeviceID:           "7",
		SubjectID:          "10",
		Bpm:                62,
		BodyTemperature:    36.5,
		AmbientTemperature: 22.0,
		Status:             repository.StatusFatigued,
	}
	mockInfluxRepo.On("WritePeriodic", mock.Anything, expectedParam).Return(nil)

	service.WritePeriodicData(mockMessage)
	mockInfluxRepo.AssertExpectations(t)
	mockDeviceRepo.AssertExpectations(t)
	mockSubjectRepo.AssertExpectations(t)
}

func TestWritePeriodicDeviceNotExists(t *testing.T) {
	mockInfluxRepo := new(mocks.MockInfluxRepository)
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	mockSubjectRepo := new(mocks.MockSubjectRepository)
	service := NewInfluxService(mockInfluxRepo, mockSubjectRepo, mockDeviceRepo)

	mockMessage := &mocks.MockMessage{
		MockPayload: []byte(`{"bpm":72,"suhu_objek":36.5,"suhu_lingkungan":22.0,"status":0}`),
		MockTopic:   "esp32/bb:cc:dd:ee:ff:gg/periodic",
	}

	mockNilDevice := (*repository.Device)(nil)
	mockDeviceRepo.On("FindDeviceByClientID", mock.Anything, "bb:cc:dd:ee:ff:gg").Return(mockNilDevice, nil)

	mockDevice := repository.Device{ClientID: "bb:cc:dd:ee:ff:gg"}
	mockResult := &mocks.MockResult{}
	mockDeviceRepo.On("AddDevice", mock.Anything, mockDevice).Return(mockResult, nil)

	mockResult.On("LastInsertId").Return(int64(8), nil)

	mockSubject := &repository.Subject{ID: 10, Name: "My Subject", IsFatigued: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockSubjectRepo.On("FindSubjectByDeviceID", mock.Anything, int64(8)).Return(mockSubject, nil)

	mockSubjectRepo.On("UpdateSubjectFatiguedStatus", mock.Anything, repository.Subject{
		ID:         10,
		IsFatigued: false,
	}).Return(mockResult, nil)

	expectedParam := repository.InfluxPeriodicPointParam{
		DeviceID:           "8",
		SubjectID:          "10",
		Bpm:                72,
		BodyTemperature:    36.5,
		AmbientTemperature: 22.0,
		Status:             repository.StatusNotFatigued,
	}
	mockInfluxRepo.On("WritePeriodic", mock.Anything, expectedParam).Return(nil)

	service.WritePeriodicData(mockMessage)
	mockInfluxRepo.AssertExpectations(t)
	mockDeviceRepo.AssertExpectations(t)
	mockSubjectRepo.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

func TestWritePerpetualDataDeviceExists(t *testing.T) {
	mockInfluxRepo := new(mocks.MockInfluxRepository)
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	mockSubjectRepo := new(mocks.MockSubjectRepository)
	service := NewInfluxService(mockInfluxRepo, mockSubjectRepo, mockDeviceRepo)

	mockMessage := &mocks.MockMessage{
		MockPayload: []byte(`{"raw_ecg":70.0}`),
		MockTopic:   "esp32/aa:bb:cc:dd:ee:ff/perpetual",
	}

	mockDevice := &repository.Device{ID: 7, ClientID: "aa:bb:cc:dd:ee:ff"}
	mockSubject := &repository.Subject{ID: 10, Name: "My Subject", IsFatigued: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mockDeviceRepo.On("FindDeviceByClientID", mock.Anything, "aa:bb:cc:dd:ee:ff").Return(mockDevice, nil)
	mockSubjectRepo.On("FindSubjectByDeviceID", mock.Anything, int64(7)).Return(mockSubject, nil)

	mockInfluxRepo.On("WritePerpetual", mock.Anything, mock.MatchedBy(func(p repository.InfluxPerpetualPointParam) bool {
		return p.DeviceID == "7" && p.SubjectID == "10" && p.RawEcg == 70.0
	})).Return(nil)

	service.WritePerpetualData(mockMessage)
	mockInfluxRepo.AssertExpectations(t)
	mockDeviceRepo.AssertExpectations(t)
	mockSubjectRepo.AssertExpectations(t)
}

func TestWritePerpetualDataDeviceNotExists(t *testing.T) {
	mockInfluxRepo := new(mocks.MockInfluxRepository)
	mockDeviceRepo := new(mocks.MockDeviceRepository)
	mockSubjectRepo := new(mocks.MockSubjectRepository)
	service := NewInfluxService(mockInfluxRepo, mockSubjectRepo, mockDeviceRepo)

	mockMessage := &mocks.MockMessage{
		MockPayload: []byte(`{"raw_ecg":80.0}`),
		MockTopic:   "esp32/bb:cc:dd:ee:ff:gg/perpetual",
	}

	mockNilDevice := (*repository.Device)(nil)
	mockDeviceRepo.On("FindDeviceByClientID", mock.Anything, "bb:cc:dd:ee:ff:gg").Return(mockNilDevice, nil)

	mockDevice := repository.Device{ClientID: "bb:cc:dd:ee:ff:gg"}
	mockResult := &mocks.MockResult{}
	mockDeviceRepo.On("AddDevice", mock.Anything, mockDevice).Return(mockResult, nil)

	mockResult.On("LastInsertId").Return(int64(8), nil)

	mockSubject := &repository.Subject{ID: 10, Name: "My Subject", IsFatigued: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockSubjectRepo.On("FindSubjectByDeviceID", mock.Anything, int64(8)).Return(mockSubject, nil)

	mockInfluxRepo.On("WritePerpetual", mock.Anything, mock.MatchedBy(func(p repository.InfluxPerpetualPointParam) bool {
		return p.DeviceID == "8" && p.SubjectID == "10" && p.RawEcg == 80.0
	})).Return(nil)

	service.WritePerpetualData(mockMessage)
	mockInfluxRepo.AssertExpectations(t)
	mockDeviceRepo.AssertExpectations(t)
	mockSubjectRepo.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}
