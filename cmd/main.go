package main

import (
	api "github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
)

func main() {
	// ロガーを初期化
	log := logger.NewLogger()

	if err := api.Run(log); err != nil {
		log.Error("Failed to start server :%v", err)
	}
}
