package dto

import (
	"caps_influx/internal/repository"
	"time"
)

type DeviceResponse struct {
	ID        int64
	ClientID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToDeviceResponses(devices *[]repository.Device) []DeviceResponse {
	responses := []DeviceResponse{}
	for _, d := range *devices {
		response := DeviceResponse{
			ID:        d.ID,
			ClientID:  d.ClientID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		}
		responses = append(responses, response)
	}
	return responses
}

type DeviceCreateReq struct {
	ClientID string `json:"client_id"`
}
