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
	if ctx == nil {
		ctx = context.Background()
	}

	container := dig.New()
	container.Provide(NewEcho)
	container.Provide(handler.NewHandler)
	return container.Invoke(func(e *echo.Echo, ha *handler.Handler) error {
		RegisterHandlers(e, ha)

		addr := ":8080"
		srv := &http.Server{
			Addr: addr,
		}

		// サーバ全体のライフサイクル（起動〜停止）を表すctx。
		// shutdown開始時に明示的にcancelすることで、起動ログ等にも一貫したctxが渡されます。
		serverCtx, cancelServer := context.WithCancel(ctx)
		defer cancelServer()

		errCh := make(chan error, 1)
		go func() {
			logger.Info(serverCtx, "server starting", "addr", addr)
			// StartServer は shutdown 経由の停止時も http.ErrServerClosed を返します。
			// それ以外のエラーは異常終了として扱います。
			err := e.StartServer(srv)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
				return
			}
			// 正常終了（shutdown）も通知しておくと、select側でgoroutineリークを避けられます。
			errCh <- nil
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigCh)

		select {
		case err := <-errCh:
			if err != nil {
				logger.Error(serverCtx, "server stopped with error", "err", err)
				return err
			}
		case sig := <-sigCh:
			logger.Info(serverCtx, "shutdown signal received", "signal", sig.String())
		case <-ctx.Done():
			logger.Info(serverCtx, "context canceled", "err", ctx.Err())
		}

		cancelServer()

		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			logger.Error(ctx, "graceful shutdown failed", "err", err)
			return err
		}
		logger.Info(ctx, "server shutdown complete")
		return nil
	})
}

func NewEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	return e
}
