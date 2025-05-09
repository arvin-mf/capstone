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
