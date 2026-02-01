package logger

import (
	"errors"
	"os"
	"strings"
	"sync"
)

// LoggerInterface defines the interface for the logger
type LoggerInterface interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}


var (
	mu      sync.RWMutex
	current LoggerInterface
)

func NewLogger() error {
	mu.Lock()
	defer mu.Unlock()

	if current != nil {
		return nil
	}

	impl := strings.ToLower(strings.TrimSpace(os.Getenv("LOGGER")))
	switch impl {
	case "zap":
		current = NewZapLogger()
	case "", "simple":
		current = NewSimpleLogger()
	default:
		return errors.New("unknown LOGGER value: " + impl)
	}

	return nil
}

func CloseLogger() {
	mu.Lock()
	defer mu.Unlock()

	if z, ok := current.(*zapLogger); ok {
		_ = z.logger.Sync()
	}
	current = nil
}

func get() LoggerInterface {
	mu.RLock()
	l := current
	mu.RUnlock()
	if l != nil {
		return l
	}
	mu.Lock()
	defer mu.Unlock()
	if current == nil {
		current = NewSimpleLogger()
	}
	return current
}

func Info(msg string, fields ...interface{})  { get().Info(msg, fields...) }
func Error(msg string, fields ...interface{}) { get().Error(msg, fields...) }
func Debug(msg string, fields ...interface{}) { get().Debug(msg, fields...) }
func Warn(msg string, fields ...interface{})  { get().Warn(msg, fields...) }
