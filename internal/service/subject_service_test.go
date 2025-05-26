package service

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/mocks"
	"caps_influx/internal/repository"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSubjects(t *testing.T) {
	mockSubjectRepo := new(mocks.MockSubjectRepository)
	service := NewSubjectService(mockSubjectRepo)

	subjects := []repository.Subject{
		{ID: 1, Name: "Subject 1", CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Subject 2", CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
	}
	mockSubjectRepo.On("FindSubjects", mock.Anything).Return(subjects, nil)

	result, err := service.GetSubjects(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, dto.ToSubjectResponses(&subjects), result)
	mockSubjectRepo.AssertExpectations(t)
}

func TestGetSubjectsEmpty(t *testing.T) {
	mockSubjectRepo := new(mocks.MockSubjectRepository)
	service := NewSubjectService(mockSubjectRepo)

	emptySubjects := []repository.Subject{}
	mockSubjectRepo.On("FindSubjects", mock.Anything).Return(emptySubjects, nil)

	result, err := service.GetSubjects(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, dto.ToSubjectResponses(&emptySubjects), result)
	mockSubjectRepo.AssertExpectations(t)
}

func TestAddSubject(t *testing.T) {
	mockSubjectRepo := new(mocks.MockSubjectRepository)
	service := NewSubjectService(mockSubjectRepo)

	subject := dto.SubjectCreateReq{Name: "New Subject"}
	mockResult := &mocks.MockResult{}
	mockSubjectRepo.On("AddSubject", mock.Anything, mock.AnythingOfType("repository.Subject")).Return(mockResult, nil)

	err := service.AddSubject(context.Background(), subject)
	assert.NoError(t, err)
	mockSubjectRepo.AssertExpectations(t)
}

func TestDeleteSubject(t *testing.T) {
	mockSubjectRepo := new(mocks.MockSubjectRepository)
	service := NewSubjectService(mockSubjectRepo)

	subjectID := int64(1)
	mockResult := &mocks.MockResult{}
	mockSubjectRepo.On("DeleteSubject", mock.Anything, mock.AnythingOfType("repository.Subject")).Return(mockResult, nil)

	err := service.DeleteSubject(context.Background(), subjectID)
	assert.NoError(t, err)
	mockSubjectRepo.AssertExpectations(t)
}
