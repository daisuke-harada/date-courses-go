package iface

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/config"
	"github.com/daisuke-harada/date-courses-go/internal/di"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Run(ctx context.Context) error {
	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container := di.NewContainer()
	container.MustProvide(NewEcho)
	container.MustProvide(config.Get)
	container.MustProvide(di.ProvideDB)
	di.BuildContainer(container)

	return container.Invoke(func(e *echo.Echo) error {
		openapi.RegisterHandlers(e, handler.NewHandler(container))

		addr := ":7777"
		srv := &http.Server{
			Addr: addr,
		}

		// errCh をバッファ1にしているのは、select が notifyCtx.Done() 側を先に選ぶ可能性があるためです。
		// もし errCh がバッファ0だと、select が notifyCtx.Done() を選んで先に進んだ後は errCh を受信する箇所が無くなるので、
		// StartServer が shutdown の結果として戻ってきたタイミングで errCh <- nil/err がブロックし、
		// 起動goroutineが終了できずに残り続ける(= goroutine leak っぽい状態)可能性があります。
		errCh := make(chan error, 1)
		go func() {
			slog.InfoContext(notifyCtx, "server starting", "addr", addr)
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
				slog.ErrorContext(notifyCtx, "server stopped with error", "err", err)
				return err
			}
		case <-notifyCtx.Done():
			// notifyCtx はシグナルでも親ctxのキャンセルでも Done になります。
			slog.InfoContext(notifyCtx, "context canceled", "err", notifyCtx.Err())
		}

		shutdownCtx, cancel := context.WithTimeout(notifyCtx, 10*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			slog.ErrorContext(ctx, "graceful shutdown failed", "err", err)
			return err
		}
		slog.InfoContext(ctx, "server shutdown complete")
		return nil
	})
}

func NewEcho() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = middleware.CustomHTTPErrorHandler
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(middleware.CORSMiddleware())
	e.Use(middleware.RequestIDMiddleware)
	e.Use(middleware.AccessLogMiddleware)
	return e
}
