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

type DeviceWithSubjectResponse struct {
	DeviceID   int64     `json:"device_id"`
	SubjectID  int64     `json:"subject_id"`
	Name       string    `json:"name"`
	IsFatigued bool      `json:"is_fatigued"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToDeviceWithSubjectResponses(subjects *[]repository.DeviceWithSubject) []DeviceWithSubjectResponse {
	responses := []DeviceWithSubjectResponse{}
	for _, s := range *subjects {
		response := DeviceWithSubjectResponse{
			DeviceID:   s.DeviceID,
			SubjectID:  s.SubjectID.Int64,
			Name:       s.Name.String,
			IsFatigued: s.IsFatigued.Bool,
			CreatedAt:  s.CreatedAt.Time,
		}
		responses = append(responses, response)
	}
	return responses
}
