package service

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type InfluxService interface {
	WriteData(m mqtt.Message)
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

func (s *influxService) WriteData(m mqtt.Message) {
	ctx := context.Background()

	var data dto.SubscribeData
	if err := json.Unmarshal(m.Payload(), &data); err != nil {
		fmt.Printf("Failed to fetch data from broker: %v\n", err)
		return
	}

	topic := m.Topic()
	topicParts := strings.Split(topic, "/")
	if len(topicParts) < 3 {
		fmt.Printf("This topic '%s' is not in a suitable format\n", topic)
		return
	}

	clientID := topicParts[1]

	device, err := s.deviceRepo.FindDeviceByClientID(ctx, clientID)
	if err != nil {
		fmt.Println("Failed to fetch device from database")
		return
	}

	subject, err := s.subjectRepo.FindSubjectByDeviceID(ctx, device.ID)
	if err != nil {
		fmt.Println("Failed to fetch subject from database")
		return
	}

	var status repository.SubjectStatus
	if data.Status == 1 {
		status = repository.StatusFatigued
	} else {
		status = repository.StatusNotFatigued
	}

	if err := s.influxRepo.Write(ctx, repository.InfluxPointParam{
		DeviceID:           strconv.Itoa(int(device.ID)),
		SubjectID:          strconv.Itoa(int(subject.ID)),
		Bpm:                data.Bpm,
		BodyTemperature:    data.BodyTemperature,
		AmbientTemperature: data.AmbientTemperature,
		Status:             status,
	}); err != nil {
		fmt.Printf("Failed to write data to InfluxDB: %v", err)
		return
	}
}
