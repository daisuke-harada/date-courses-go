package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

// AccessLogMiddleware はリクエストごとに slog でアクセスログを出力するミドルウェア。
// pkg/logger の contextAttrsHandler を通じて request_id が自動的にログに含まれる。
func AccessLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		req := c.Request()
		res := c.Response()
		latency := time.Since(start)

		// リクエストの context には apimw.RequestIDMiddleware により request_id が
		// 注入されているため、slog.InfoContext に渡すだけで自動的にログに含まれる。
		slog.InfoContext(
			req.Context(),
			"access",
			slog.String("method", req.Method),
			slog.String("uri", req.RequestURI),
			slog.Int("status", res.Status),
			slog.String("latency", latency.String()),
			slog.String("remote_ip", c.RealIP()),
		)

		return err
	}
}
