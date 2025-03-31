package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

func (m *Middleware) LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start).Milliseconds()

			// Warn: 処理時間が長い場合に警告を記録
			if latency > 1000 { // 例: 1秒を超える場合
				m.log.Warn("Request took too long",
					"method", c.Request().Method,
					"path", c.Request().URL.Path,
					"status", c.Response().Status,
					"latency", latency,
				)
			}
			// リクエストの詳細をログに記録
			m.log.Info("Request",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency", latency,
			)

			// エラーが発生した場合はログに記録
			if err != nil {
				c.Error(err) // Echoのデフォルトエラーハンドラーを呼び出す
				m.log.Error("Request error",
					"method", c.Request().Method,
					"path", c.Request().URL.Path,
					"error", err,
				)
			}

			return err
		}
	}
}
