package main

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
)

func main() {
	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Close()

	if err := api.Run(context.Background()); err != nil {
		logger.Error(context.Background(), "Failed to start server", "err", err)
	}
}
