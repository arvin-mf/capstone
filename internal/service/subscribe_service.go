package service

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SubscribeService interface {
	SubscribeData()
}

type subscribeService struct {
	client     mqtt.Client
	influxServ InfluxService
}

func NewSubscribeService(c mqtt.Client, is InfluxService) SubscribeService {
	return &subscribeService{
		client:     c,
		influxServ: is,
	}
}

func (s *subscribeService) SubscribeData() {
	s.client.Subscribe("esp32/+/discrete", 1, func(c mqtt.Client, m mqtt.Message) {
		s.influxServ.WriteData(m)
	})
}
