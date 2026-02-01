package logger

import (
	"fmt"
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
	l.logger.Print("INFO: " + fmt.Sprintf(msg, fields...))
}

func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
	l.logger.Print("ERROR: " + fmt.Sprintf(msg, fields...))
}

func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Print("DEBUG: " + fmt.Sprintf(msg, fields...))
}

func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Print("WARN: " + fmt.Sprintf(msg, fields...))
}
