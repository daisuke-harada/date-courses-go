package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newEchoWithLoginRateLimit(limit int) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = middleware.CustomHTTPErrorHandler
	e.Use(middleware.LoginRateLimitMiddleware(limit))
	e.POST("/api/v1/login", dummyHandler)
	e.POST("/api/v1/signup", dummyHandler)
	e.GET("/api/v1/date_spots", dummyHandler)
	return e
}

func post(t *testing.T, e *echo.Echo, path, ip string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, nil)
	req.RemoteAddr = ip + ":12345"
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code
}

func TestLoginRateLimitMiddleware(t *testing.T) {
	// ブルートフォースを防ぐため、同一 IP からのログイン試行に上限を設ける
	t.Run("rejects_after_limit_exceeded", func(t *testing.T) {
		e := newEchoWithLoginRateLimit(3)

		for i := 0; i < 3; i++ {
			assert.Equal(t, http.StatusOK, post(t, e, "/api/v1/login", "10.0.0.1"), "%d 回目は通るはず", i+1)
		}
		assert.Equal(t, http.StatusTooManyRequests, post(t, e, "/api/v1/login", "10.0.0.1"), "上限を超えたら 429")
	})

	// 制限は IP ごと。別の利用者を巻き込まない
	t.Run("limits_per_ip", func(t *testing.T) {
		e := newEchoWithLoginRateLimit(2)

		post(t, e, "/api/v1/login", "10.0.0.1")
		post(t, e, "/api/v1/login", "10.0.0.1")
		assert.Equal(t, http.StatusTooManyRequests, post(t, e, "/api/v1/login", "10.0.0.1"))
		assert.Equal(t, http.StatusOK, post(t, e, "/api/v1/login", "10.0.0.2"), "別 IP は影響を受けない")
	})

	// 新規登録もアカウント作成の乱用対象になるため対象に含める
	t.Run("applies_to_signup", func(t *testing.T) {
		e := newEchoWithLoginRateLimit(1)

		assert.Equal(t, http.StatusOK, post(t, e, "/api/v1/signup", "10.0.0.3"))
		assert.Equal(t, http.StatusTooManyRequests, post(t, e, "/api/v1/signup", "10.0.0.3"))
	})

	// 通常の閲覧は制限しない
	t.Run("does_not_limit_other_routes", func(t *testing.T) {
		e := newEchoWithLoginRateLimit(1)

		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/date_spots", nil)
			req.RemoteAddr = "10.0.0.4:12345"
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}
