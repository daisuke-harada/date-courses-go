//go:generate oapi-codegen -generate types -o api_types.gen.go -package api ../../../../api/resolved/openapi/openapi.yaml
//go:generate oapi-codegen -generate echo-server -o api_server.gen.go -package api ../../../../api/resolved/openapi/openapi.yaml
//go:generate go run handler_generator.go

package api

import (
	"context"
	"errors"
	"net/http"
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
	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container := dig.New()
	container.Provide(NewEcho)
	container.Provide(handler.NewHandler)

	return container.Invoke(func(e *echo.Echo, handler *handler.Handler) error {
		RegisterHandlers(e, handler)

		addr := ":8080"
		srv := &http.Server{
			Addr: addr,
		}

		// errCh をバッファ1にしているのは、select が ctx.Done() 側を先に選ぶ可能性があるためです。
		// もし errCh がバッファ0だと、select が ctx.Done() を選んで先に進んだ後は errCh を受信する箇所が無くなるので、
		// StartServer が shutdown の結果として戻ってきたタイミングで errCh <- nil/err がブロックし、
		// 起動goroutineが終了できずに残り続ける(= goroutine leak っぽい状態)可能性があります。
		errCh := make(chan error, 1)
		go func() {
			logger.Info(notifyCtx, "server starting", "addr", addr)
			// StartServer は shutdown 経由の停止時も http.ErrServerClosed を返します。
			// それ以外のエラーは異常終了として扱います。
			err := e.StartServer(srv)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
				return
			}
			errCh <- nil
		}()

		select {
		case err := <-errCh:
			if err != nil {
				logger.Error(notifyCtx, "server stopped with error", "err", err)
				return err
			}
		case <-notifyCtx.Done():
			// notifyCtx はシグナルでも親ctxのキャンセルでも Done になります。
			logger.Info(notifyCtx, "context canceled", "err", notifyCtx.Err())
		}

		shutdownCtx, cancel := context.WithTimeout(notifyCtx, 10*time.Second)
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
