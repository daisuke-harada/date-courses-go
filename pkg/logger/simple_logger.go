package logger

import (
	"log"
	"os"
)

type SimpleLogger struct {
	logger *log.Logger
}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	l.logger.Printf("INFO: "+msg, fields...)
}

func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
	l.logger.Printf("ERROR: "+msg, fields...)
}

func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Printf("DEBUG: "+msg, fields...)
}

func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Printf("WARN: "+msg, fields...)
}
