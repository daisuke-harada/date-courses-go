//go:generate oapi-codegen -generate types -o api_types.gen.go -package api ../../../../api/resolved/openapi/openapi.yaml
//go:generate oapi-codegen -generate echo-server -o api_server.gen.go -package api ../../../../api/resolved/openapi/openapi.yaml
//go:generate go run handler_generator.go

package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api/handler"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

func Run(ctx context.Context) error {
	container := dig.New()
	container.Provide(NewEcho)
	container.Provide(handler.NewHandler)
	return container.Invoke(func(e *echo.Echo, ha *handler.Handler) error {
		RegisterHandlers(e, ha)

		addr := ":8080"
		srv := &http.Server{
			Addr:    addr,
			Handler: e,
		}

		errCh := make(chan error, 1)
		go func() {
			logger.Info("server starting on %s", addr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigCh)

		select {
		case err := <-errCh:
			logger.Error("server stopped with error: %v", err)
			return err
		case sig := <-sigCh:
			logger.Info("shutdown signal received: %s", sig.String())
		case <-ctx.Done():
			logger.Info("context canceled: %v", ctx.Err())
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed: %v", err)
			return err
		}
		logger.Info("server shutdown complete")
		return nil
	})
}

func NewEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	return e
}
