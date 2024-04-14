package mogs

import (
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/stretchr/testify/mock"
)

type mockLogger struct {
	mock.Mock
}

func NewMockLogger() *mockLogger {
	return &mockLogger{}
}

func (m *mockLogger) Fatal(err error, fields ...logger.LogFields) {
}

func (m *mockLogger) Error(err error, fields ...logger.LogFields) {
}

func (m *mockLogger) Warn(err error, fields ...logger.LogFields) {
}

func (m *mockLogger) Info(message string, fields ...logger.LogFields) {
}

func (m *mockLogger) Debug(message string, fields ...logger.LogFields) {
}
