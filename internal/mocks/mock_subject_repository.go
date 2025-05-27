package mocks

import (
	"caps_influx/internal/repository"
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockSubjectRepository struct {
	mock.Mock
}

func (m *MockSubjectRepository) FindSubjects(ctx context.Context) ([]repository.Subject, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repository.Subject), args.Error(1)
}

func (m *MockSubjectRepository) AddSubject(ctx context.Context, params repository.Subject) (sql.Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSubjectRepository) DeleteSubject(ctx context.Context, params repository.Subject) (sql.Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockSubjectRepository) FindSubjectByDeviceID(ctx context.Context, dID int64) (*repository.Subject, error) {
	args := m.Called(ctx, dID)
	return args.Get(0).(*repository.Subject), args.Error(1)
}

func (m *MockSubjectRepository) UpdateSubjectFatiguedStatus(ctx context.Context, params repository.Subject) (sql.Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sql.Result), args.Error(1)
}
