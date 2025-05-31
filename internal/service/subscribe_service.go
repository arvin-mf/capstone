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

type SubscribeService interface {
	SubscribePeriodicData()
	SubscribePerpetualData()
	SubscribeStatusData()
}

type subscribeService struct {
	client     mqtt.Client
	influxServ InfluxService
	deviceRepo repository.DeviceRepository
}

func NewSubscribeService(c mqtt.Client, is InfluxService, dr repository.DeviceRepository) SubscribeService {
	return &subscribeService{
		client:     c,
		influxServ: is,
		deviceRepo: dr,
	}
}

func (s *subscribeService) SubscribePeriodicData() {
	s.client.Subscribe("esp32/+/discrete", 1, func(c mqtt.Client, m mqtt.Message) {
		s.influxServ.WritePeriodicData(m)
	})
}

func (s *subscribeService) SubscribePerpetualData() {
	s.client.Subscribe("esp32/+/continue", 1, func(c mqtt.Client, m mqtt.Message) {
		s.influxServ.WritePerpetualData(m)
	})
}

func (s *subscribeService) SubscribeStatusData() {
	s.client.Subscribe("esp32/+/status", 1, func(c mqtt.Client, m mqtt.Message) {
		s.updateDeviceStatus(m)
	})
}

func (s *subscribeService) updateDeviceStatus(m mqtt.Message) {
	topic := m.Topic()
	topicParts := strings.Split(topic, "/")
	if len(topicParts) < 3 {
		fmt.Println("this topic is not in a suitable format")
		return
	}
	clientID := topicParts[1]

	ctx := context.Background()

	var req dto.SubscribeStatusData
	if err := json.Unmarshal(m.Payload(), &req); err != nil {
		fmt.Printf("Failed to unmarshal message payload: %v\n", err)
		return
	}

	device, err := s.deviceRepo.FindDeviceByClientID(ctx, clientID)
	if err != nil {
		fmt.Printf("Failed to find device by client ID %s: %v\n", clientID, err)
		return
	}

	if device != nil {
		if req.DeviceStatus == "on" && !device.DeviceStatus {
			_, err = s.deviceRepo.UpdateDeviceStatus(ctx, repository.Device{
				ID:           device.ID,
				DeviceStatus: true,
			})
			if err != nil {
				fmt.Printf("Failed to update device status: %v\n", err)
				return
			}
		}

		key := "device-" + strconv.Itoa(int(device.ID))

		err = s.deviceRepo.SetDeviceStatusToRedis(ctx, key, "on")
		if err != nil {
			fmt.Printf("Failed to set device status to Redis: %v\n", err)
			return
		}
	}
}
