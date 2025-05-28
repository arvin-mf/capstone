package mocks

import (
	"caps_influx/internal/repository"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockInfluxRepository struct {
	mock.Mock
}

func (m *MockInfluxRepository) WritePeriodic(ctx context.Context, params repository.InfluxPeriodicPointParam) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockInfluxRepository) WritePerpetual(ctx context.Context, params repository.InfluxPerpetualPointParam) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}
