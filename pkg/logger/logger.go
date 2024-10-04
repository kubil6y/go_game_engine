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
	LEVEL_DEBUG LogLevel = iota
	LEVEL_INFO
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_OFF
)

func (l LogLevel) String() string {
	switch l {
	case LEVEL_DEBUG:
		return "DEBUG"
	case LEVEL_INFO:
		return "INFO"
	case LEVEL_ERROR:
		return "ERROR"
	case LEVEL_FATAL:
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
	case LEVEL_DEBUG:
		return Green
	case LEVEL_INFO:
		return Blue
	case LEVEL_ERROR:
		return Red
	case LEVEL_FATAL:
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
		minLevel: LEVEL_INFO,
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
	l.print(LEVEL_DEBUG, message, properties)
}

func (l *Logger) Info(message string, properties map[string]any) {
	l.print(LEVEL_INFO, message, properties)
}

func (l *Logger) Error(err error, message string, properties map[string]any) {
	wrapped := fmt.Errorf("%s: %w", message, err)
	l.print(LEVEL_ERROR, wrapped.Error(), properties)
}

func (l *Logger) Fatal(err error, message string, properties map[string]any) {
	wrapped := fmt.Errorf("%s: %w", message, err)
	l.print(LEVEL_FATAL, wrapped.Error(), properties)
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

	if level >= LEVEL_ERROR {
		logMessage += "\n" + string(debug.Stack())
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.out, logMessage)
}
