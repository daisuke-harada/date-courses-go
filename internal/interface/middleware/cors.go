package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// CORSMiddleware はCORSの設定を行うミドルウェアです。
// 許可するオリジンは環境変数 CORS_ALLOW_ORIGINS で指定します。
// ハードコードすると本番のフロントエンドから叩けなくなるため、設定から受け取ります。
func CORSMiddleware(allowOrigins []string) echo.MiddlewareFunc {
	return echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Accept", "Content-Type", "Authorization"},
	})
}
