package unit

import (
	"bytes"
	"strings"
	"testing"

	"memwright/api/pkg/logger"
)

func TestLoggerInfo(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo)

	log.Info("test message %s", "hello")

	output := buf.String()
	if !strings.Contains(output, "[INFO ]") {
		t.Errorf("expected [INFO ], got: %s", output)
	}
	if !strings.Contains(output, "test message hello") {
		t.Errorf("expected 'test message hello', got: %s", output)
	}
}

func TestLoggerWarn(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo)

	log.Warn("warning message")

	if !strings.Contains(buf.String(), "[WARN ]") {
		t.Errorf("expected [WARN ], got: %s", buf.String())
	}
}

func TestLoggerError(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo)

	log.Error("error message")

	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Errorf("expected [ERROR], got: %s", buf.String())
	}
}

func TestLoggerDebugEnabled(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelDebug)

	log.Debug("debug message")

	if !strings.Contains(buf.String(), "[DEBUG]") {
		t.Errorf("expected [DEBUG], got: %s", buf.String())
	}
}

func TestLoggerDebugFiltered(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo)

	log.Debug("debug message")

	if buf.String() != "" {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestLoggerMinLevel(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelWarn)

	log.Info("info message")
	log.Warn("warning message")

	output := buf.String()
	if strings.Contains(output, "info message") {
		t.Error("info should be filtered")
	}
	if !strings.Contains(output, "warning message") {
		t.Error("warning should be logged")
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    logger.Level
		expected string
	}{
		{logger.LevelDebug, "DEBUG"},
		{logger.LevelInfo, "INFO"},
		{logger.LevelWarn, "WARN"},
		{logger.LevelError, "ERROR"},
	}

	for _, tt := range tests {
		if tt.level.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.level.String())
		}
	}
}

func TestLoggerImplementsInterface(t *testing.T) {
	var _ logger.Logger = logger.New(nil, logger.LevelInfo)
}
