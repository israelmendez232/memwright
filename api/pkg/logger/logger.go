package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger is the interface for dependency injection.
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

// StandardLogger implements Logger.
type StandardLogger struct {
	mu       sync.Mutex
	out      io.Writer
	minLevel Level
}

// New creates a logger. Defaults to stdout and Info level.
func New(out io.Writer, minLevel Level) *StandardLogger {
	if out == nil {
		out = os.Stdout
	}
	return &StandardLogger{out: out, minLevel: minLevel}
}

func (l *StandardLogger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

func (l *StandardLogger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

func (l *StandardLogger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

func (l *StandardLogger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

func (l *StandardLogger) Fatal(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
	os.Exit(1)
}

func (l *StandardLogger) log(level Level, format string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	ts := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.out, "[%s] [%-5s] %s\n", ts, level, msg)
}
