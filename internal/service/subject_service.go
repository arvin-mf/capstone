package service

import (
	"caps_influx/internal/dto"
	"caps_influx/internal/repository"
	"context"
)

type SubjectService interface {
	GetSubjects(ctx context.Context) ([]dto.SubjectResponse, error)
	AddSubject(ctx context.Context, arg dto.SubjectCreateReq) error
	DeleteSubject(ctx context.Context, subjectID int64) error
}

type subjectService struct {
	subjectRepo repository.SubjectRepository
}

func NewSubjectService(sr repository.SubjectRepository) SubjectService {
	return &subjectService{
		subjectRepo: sr,
	}
}

func (s *subjectService) GetSubjects(ctx context.Context) ([]dto.SubjectResponse, error) {
	subjects, err := s.subjectRepo.FindSubjects(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToSubjectResponses(&subjects), nil
}

func (s *subjectService) AddSubject(ctx context.Context, arg dto.SubjectCreateReq) error {
	_, err := s.subjectRepo.AddSubject(ctx, repository.Subject{
		Name: arg.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *subjectService) DeleteSubject(ctx context.Context, subjectID int64) error {
	_, err := s.subjectRepo.DeleteSubject(ctx, repository.Subject{
		ID: subjectID,
	})
	if err != nil {
		return err
	}

	return nil
}
