package service

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type InfluxService interface {
	WritePeriodicData(m mqtt.Message)
	WritePerpetualData(m mqtt.Message)
}

type influxService struct {
	influxRepo  repository.InfluxRepository
	deviceRepo  repository.DeviceRepository
	subjectRepo repository.SubjectRepository
}

func NewInfluxService(ir repository.InfluxRepository, sr repository.SubjectRepository, dr repository.DeviceRepository) InfluxService {
	return &influxService{
		influxRepo:  ir,
		subjectRepo: sr,
		deviceRepo:  dr,
	}
}

func (s *influxService) WritePeriodicData(m mqtt.Message) {
	ctx := context.Background()

	var data dto.SubscribePeriodicData
	if err := json.Unmarshal(m.Payload(), &data); err != nil {
		fmt.Printf("Failed to fetch data from broker: %v\n", err)
		return
	}

	deviceID, subjectID, err := s.findDeviceAndSubject(ctx, m)
	if err != nil {
		fmt.Printf("Failed to fetch data from database: %v", err)
		return
	}

	var (
		isFatigued  bool
		pointStatus repository.SubjectStatus
	)

	if data.Status == 1 {
		pointStatus = repository.StatusFatigued
		isFatigued = true
	} else {
		pointStatus = repository.StatusNotFatigued
		isFatigued = false
	}

	subjectIDInt, err := strconv.Atoi(subjectID)
	if err != nil {
		fmt.Printf("Failed to parse subject ID: %v", err)
		return
	}

	_, err = s.subjectRepo.UpdateSubjectFatiguedStatus(ctx, repository.Subject{
		ID:         int64(subjectIDInt),
		IsFatigued: isFatigued,
	})
	if err != nil {
		fmt.Printf("Failed to update fatigued status: %v", err)
		return
	}

	err = s.influxRepo.WritePeriodic(ctx, repository.InfluxPeriodicPointParam{
		DeviceID:           deviceID,
		SubjectID:          subjectID,
		Bpm:                data.Bpm,
		BodyTemperature:    data.BodyTemperature,
		AmbientTemperature: data.AmbientTemperature,
		Status:             pointStatus,
	})
	if err != nil {
		fmt.Printf("Failed to write data to InfluxDB: %v", err)
		return
	}
}

func (s *influxService) WritePerpetualData(m mqtt.Message) {
	ctx := context.Background()

	var data dto.SubscribePerpetualData
	if err := json.Unmarshal(m.Payload(), &data); err != nil {
		fmt.Printf("Failed to fetch data from broker: %v", err)
		return
	}

	deviceID, subjectID, err := s.findDeviceAndSubject(ctx, m)
	if err != nil {
		fmt.Printf("Failed to fetch data from database")
		return
	}

	err = s.influxRepo.WritePerpetual(ctx, repository.InfluxPerpetualPointParam{
		DeviceID:  deviceID,
		SubjectID: subjectID,
		RawEcg:    data.RawEcg,
		Timestamp: time.Now().UnixNano(),
	})
	if err != nil {
		fmt.Printf("Failed to write data to InfluxDB: %v", err)
		return
	}
}

func (s *influxService) findDeviceAndSubject(ctx context.Context, m mqtt.Message) (string, string, error) {
	topic := m.Topic()
	topicParts := strings.Split(topic, "/")
	if len(topicParts) < 3 {
		return "", "", errors.New("this topic is not in a suitable format")
	}

	clientID := topicParts[1]

	device, err := s.deviceRepo.FindDeviceByClientID(ctx, clientID)
	if err != nil {
		return "", "", err
	}

	subject, err := s.subjectRepo.FindSubjectByDeviceID(ctx, device.ID)
	if err != nil {
		return "", "", err
	}

	return strconv.Itoa(int(device.ID)), strconv.Itoa(int(subject.ID)), nil
}
