package mocks

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockResult struct {
	mock.Mock
}

var _ sql.Result = &MockResult{}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockResult) RowsAffected() (int64, error) { return 1, nil }

type MockMessage struct {
	MockPayload []byte
	MockTopic   string
}

func (m MockMessage) Duplicate() bool   { return false }
func (m MockMessage) Qos() byte         { return 0 }
func (m MockMessage) Retained() bool    { return false }
func (m MockMessage) Topic() string     { return m.MockTopic }
func (m MockMessage) MessageID() uint16 { return 0 }
func (m MockMessage) Payload() []byte   { return m.MockPayload }
func (m MockMessage) Ack()              {}
