package main

import (
	"context"

	api "github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
)

func main() {
	// ロガーを初期化
	if err := logger.NewLogger(); err != nil {
		panic(err)
	}
	defer logger.CloseLogger()

	if err := api.Run(context.Background()); err != nil {
		logger.Error("Failed to start server :%v", err)
	}
}
