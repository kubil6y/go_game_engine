package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelError
	LevelFatal
	LevelOff
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		panic("Invalid log level")
	}
}

type LoggerOption func(l *Logger) *Logger

type Logger struct {
	out      io.Writer
	minLevel LogLevel
	mu       sync.Mutex
}

func New(opts ...LoggerOption) *Logger {
	l := &Logger{
		out:      os.Stdout,
		minLevel: LevelInfo,
		mu:       sync.Mutex{},
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func WithLogLevel(minLevel LogLevel) LoggerOption {
	return func(l *Logger) *Logger {
		l.minLevel = minLevel
		return l
	}
}

func (l *Logger) Debug(message string, properties map[string]any) {
	l.print(LevelDebug, message, properties)
}

func (l *Logger) Info(message string, properties map[string]any) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Error(err error, message string, properties map[string]any) {
	wrappedError := fmt.Errorf("%s: %w", message, err)
	l.print(LevelError, wrappedError.Error(), properties)
}

func (l *Logger) Fatal(err error, message string, properties map[string]any) {
	wrappedError := fmt.Errorf("%s: %w", message, err)
	l.print(LevelFatal, wrappedError.Error(), properties)
}

func (l *Logger) print(level LogLevel, message string, properties map[string]any) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string         `json:"level"`
		Time       string         `json:"time"`
		Message    string         `json:"message"`
		Properties map[string]any `json:"properties,omitempty"`
		Trace      string         `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}
	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out.Write(append(line, '\n'))
}
