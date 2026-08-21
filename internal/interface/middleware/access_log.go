package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/labstack/echo/v4"
)

// AccessLogMiddleware はリクエストごとに slog でアクセスログを出力するミドルウェア。
// pkg/logger の contextAttrsHandler を通じて request_id が自動的にログに含まれる。
func AccessLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		req := c.Request()
		latency := time.Since(start)

		// リクエストの context には apimw.RequestIDMiddleware により request_id が
		// 注入されているため、slog.InfoContext に渡すだけで自動的にログに含まれる。
		slog.InfoContext(
			req.Context(),
			"access",
			slog.String("method", req.Method),
			slog.String("uri", req.RequestURI),
			slog.Int("status", responseStatus(c, err)),
			slog.String("latency", latency.String()),
			slog.String("remote_ip", c.RealIP()),
		)

		return err
	}
}

// responseStatus は実際にクライアントへ返されるステータスコードを返します。
//
// ハンドラがエラーを返した場合、レスポンスへの書き込みは Echo の
// HTTPErrorHandler が行い、それはこのミドルウェアより後になる。
// そのため c.Response().Status をそのまま読むと、エラー時でも初期値の 200 が
// 記録されてしまう。エラーからステータスを解決することでそれを防ぐ。
func responseStatus(c echo.Context, err error) int {
	if err == nil {
		return c.Response().Status
	}

	if statusCode, _, _, ok := apperror.HTTPStatus(err); ok {
		return statusCode
	}

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Code
	}

	// 型の分からないエラーは CustomHTTPErrorHandler が 500 として扱う
	return http.StatusInternalServerError
}
