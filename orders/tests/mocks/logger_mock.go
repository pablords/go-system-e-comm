package mocks

import (
	"context"
	"log/slog"
)

// MockLogger is a mock implementation of slog.Logger for testing
type MockLogger struct {
	logs []LogEntry
}

type LogEntry struct {
	Level   slog.Level
	Message string
	Attrs   map[string]interface{}
}

func NewMockLogger() *slog.Logger {
	mock := &MockLogger{
		logs: make([]LogEntry, 0),
	}

	handler := &MockHandler{mock: mock}
	return slog.New(handler)
}

type MockHandler struct {
	mock *MockLogger
}

func (h *MockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *MockHandler) Handle(_ context.Context, r slog.Record) error {
	attrs := make(map[string]interface{})
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	h.mock.logs = append(h.mock.logs, LogEntry{
		Level:   r.Level,
		Message: r.Message,
		Attrs:   attrs,
	})

	return nil
}

func (h *MockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *MockHandler) WithGroup(name string) slog.Handler {
	return h
}

// GetLogs returns all logged entries
func (m *MockLogger) GetLogs() []LogEntry {
	return m.logs
}

// Clear removes all logged entries
func (m *MockLogger) Clear() {
	m.logs = make([]LogEntry, 0)
}
