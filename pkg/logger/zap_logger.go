package logger

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

// zapLoggerは、Zapライブラリを使用したロガーの実装です
type zapLogger struct {
	logger *zap.Logger
}

// Infoメソッドは、情報レベルのログメッセージを出力します
func (z *zapLogger) Info(msg string, fields ...interface{}) {
	z.logger.Sugar().Info(fmt.Sprintf(msg, fields...))
}

// Errorメソッドは、エラーレベルのログメッセージを出力します
func (z *zapLogger) Error(msg string, fields ...interface{}) {
	z.logger.Sugar().Error(fmt.Sprintf(msg, fields...))
}

// Debugメソッドは、デバッグレベルのログメッセージを出力します
func (z *zapLogger) Debug(msg string, fields ...interface{}) {
	z.logger.Sugar().Debug(fmt.Sprintf(msg, fields...))
}

// Warnメソッドは、警告レベルのログメッセージを出力します
func (z *zapLogger) Warn(msg string, fields ...interface{}) {
	z.logger.Sugar().Warn(fmt.Sprintf(msg, fields...))
}

// NewZapLogger関数は、新しいZapロガーを初期化し、zapLoggerとして返します
func NewZapLogger() *zapLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return &zapLogger{logger: logger}
}
