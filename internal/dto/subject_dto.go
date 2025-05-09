package dto

import (
	"caps_influx/internal/repository"
	"time"
)

type SubjectResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	IsFatigued bool      `json:"is_fatigued"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToSubjectResponses(subjects *[]repository.Subject) []SubjectResponse {
	responses := []SubjectResponse{}
	for _, s := range *subjects {
		response := SubjectResponse{
			ID:         s.ID,
			Name:       s.Name,
			IsFatigued: s.IsFatigued,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
		}
		responses = append(responses, response)
	}
	return responses
}

type SubjectCreateReq struct {
	Name string `json:"name"`
}

type SubjectWithDeviceResponse struct {
	SubjectID  int64     `json:"subject_id"`
	DeviceID   int64     `json:"device_id"`
	Name       string    `json:"name"`
	IsFatigued bool      `json:"is_fatigued"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToSubjectWithDeviceResponses(subjects *[]repository.SubjectWithDevice) []SubjectWithDeviceResponse {
	responses := []SubjectWithDeviceResponse{}
	for _, s := range *subjects {
		response := SubjectWithDeviceResponse{
			SubjectID:  s.SubjectID,
			DeviceID:   s.DeviceID,
			Name:       s.Name,
			IsFatigued: s.IsFatigued,
			CreatedAt:  s.CreatedAt,
		}
		responses = append(responses, response)
	}
	return responses
}
