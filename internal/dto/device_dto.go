package dto

import (
	"caps_influx/internal/repository"
	"time"
)

type DeviceResponse struct {
	ID           int64     `json:"id"`
	ClientID     string    `json:"client_id"`
	DeviceStatus string    `json:"device_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToDeviceResponses(devices *[]repository.Device) []DeviceResponse {
	responses := []DeviceResponse{}
	for _, d := range *devices {
		response := DeviceResponse{
			ID:       d.ID,
			ClientID: d.ClientID,
			DeviceStatus: func(status bool) string {
				if status {
					return "on"
				}
				return "off"
			}(d.DeviceStatus),
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
	DeviceID     int64     `json:"device_id"`
	DeviceStatus string    `json:"device_status"`
	UpdatedAt    time.Time `json:"updated_at"`
	SubjectID    int64     `json:"subject_id"`
	Name         string    `json:"name"`
	IsFatigued   bool      `json:"is_fatigued"`
	CreatedAt    time.Time `json:"created_at"`
}

func ToDeviceWithSubjectResponses(subjects *[]repository.DeviceWithSubject) []DeviceWithSubjectResponse {
	responses := []DeviceWithSubjectResponse{}
	for _, s := range *subjects {
		response := DeviceWithSubjectResponse{
			DeviceID: s.DeviceID,
			DeviceStatus: func(status bool) string {
				if status {
					return "on"
				}
				return "off"
			}(s.DeviceStatus),
			UpdatedAt:  s.UpdatedAt,
			SubjectID:  s.SubjectID.Int64,
			Name:       s.Name.String,
			IsFatigued: s.IsFatigued.Bool,
			CreatedAt:  s.CreatedAt.Time,
		}
		responses = append(responses, response)
	}
	return responses
}

type SetDeviceSubjectReq struct {
	SubjectID int64 `json:"subject_id"`
}
