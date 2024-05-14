package mogs

import (
	"github.com/stretchr/testify/mock"
)

type mockLogger struct {
	mock.Mock
}

func NewMockLogger() *mockLogger {
	return &mockLogger{}
}

func (m *mockLogger) Fatal(err error, fields ...map[string]any) {
}

func (m *mockLogger) Error(err error, fields ...map[string]any) {
}

func (m *mockLogger) Warn(err error, fields ...map[string]any) {
}

func (m *mockLogger) Info(message string, fields ...map[string]any) {
}

func (m *mockLogger) Debug(message string, fields ...map[string]any) {
}
