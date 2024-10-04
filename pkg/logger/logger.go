package logger

import (
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

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func getLevelColor(level LogLevel) string {
	switch level {
	case LevelDebug:
		return Green
	case LevelInfo:
		return Blue
	case LevelError:
		return Red
	case LevelFatal:
		return Red
	default:
		return Reset
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
	wrapped := fmt.Errorf("%s: %w", message, err)
	l.print(LevelError, wrapped.Error(), properties)
}

func (l *Logger) Fatal(err error, message string, properties map[string]any) {
	wrapped := fmt.Errorf("%s: %w", message, err)
	l.print(LevelFatal, wrapped.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level LogLevel, message string, properties map[string]any) {
	if level < l.minLevel {
		return
	}

	logMessage := fmt.Sprintf("%s[%s]\t%s: %s%s", getLevelColor(level), level.String(), time.Now().Format(time.RFC3339), message, Reset)

	if properties != nil {
		logMessage += " | Properties: " + fmt.Sprintf("%+v", properties)
	}

	if level >= LevelError {
		logMessage += "\n" + string(debug.Stack())
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.out, logMessage)
}
