package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
)

func main() {
	logger.Init("date-courses-go", false)
	defer logger.Close()

	if err := api.Run(context.Background()); err != nil {
		// Use slog's package-level helper (configured by logger.Init)
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}
