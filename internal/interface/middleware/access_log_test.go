package middleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureAccessLog は1リクエスト分のアクセスログを JSON として取り出します。
func captureAccessLog(t *testing.T, handler echo.HandlerFunc) map[string]any {
	t.Helper()

	var buf bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))
	defer slog.SetDefault(original)

	e := echo.New()
	e.HTTPErrorHandler = middleware.CustomHTTPErrorHandler
	e.Use(middleware.AccessLogMiddleware)
	e.GET("/test", handler)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", nil))

	var entry map[string]any
	for _, line := range bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n")) {
		var m map[string]any
		if err := json.Unmarshal(line, &m); err != nil {
			continue
		}
		if m["msg"] == "access" {
			entry = m
		}
	}
	require.NotNil(t, entry, "アクセスログが出力されていない")
	return entry
}

func TestAccessLogMiddleware(t *testing.T) {
	t.Run("logs_200_for_success", func(t *testing.T) {
		entry := captureAccessLog(t, func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		assert.Equal(t, float64(http.StatusOK), entry["status"])
	})

	// ハンドラがエラーを返した場合、レスポンスへの書き込みはミドルウェアより後に
	// 行われるため、そのまま res.Status を読むと 200 が記録されてしまう
	t.Run("logs_actual_status_when_handler_returns_error", func(t *testing.T) {
		entry := captureAccessLog(t, func(c echo.Context) error {
			return apperror.NotFound()
		})

		assert.Equal(t, float64(http.StatusNotFound), entry["status"])
	})

	t.Run("logs_500_for_internal_error", func(t *testing.T) {
		entry := captureAccessLog(t, func(c echo.Context) error {
			return apperror.InternalServerError(nil)
		})

		assert.Equal(t, float64(http.StatusInternalServerError), entry["status"])
	})

	// apperror 以外（echo が返す 405 など）も正しく記録する
	t.Run("logs_status_from_echo_http_error", func(t *testing.T) {
		entry := captureAccessLog(t, func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusMethodNotAllowed)
		})

		assert.Equal(t, float64(http.StatusMethodNotAllowed), entry["status"])
	})

	// 型の分からないエラーは 500 として扱う
	t.Run("logs_500_for_unknown_error", func(t *testing.T) {
		entry := captureAccessLog(t, func(c echo.Context) error {
			return assert.AnError
		})

		assert.Equal(t, float64(http.StatusInternalServerError), entry["status"])
	})
}
