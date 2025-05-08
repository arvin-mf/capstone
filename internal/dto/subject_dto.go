package dto

import (
	"caps_influx/internal/repository"
	"time"
)

type SubjectResponse struct {
	ID         int64
	Name       string
	IsFatigued bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
