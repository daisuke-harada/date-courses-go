package logger

// LoggerInterface defines the interface for the logger
type LoggerInterface interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// NewLogger initializes a new Zap logger and returns it as a LoggerInterface
func NewLogger() LoggerInterface {
	return NewZapLogger()
}
