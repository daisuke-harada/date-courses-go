package middleware

import (
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	iface_openapi "github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	jwtpkg "github.com/daisuke-harada/date-courses-go/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
)

// JWTAuthMiddleware は JWT Bearer トークンを検証し、currentUser をコンテキストにセットします。
// 認証が必要かどうかは iface_openapi.RequiresBearerAuth を通じて判定します。
// この関数は api/resolved/openapi/openapi.yaml の security 定義から自動生成された
// auth_routes.gen.go のマップを参照します。`make gen` で再生成されます。
func JWTAuthMiddleware(secretKey string, userRepo repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			if !iface_openapi.RequiresBearerAuth(req.Method, ctx.Path()) {
				return next(ctx)
			}

			authHeader := req.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return apperror.Unauthorized("認証が必要です。")
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := jwtpkg.Decode(tokenStr, secretKey)
			if err != nil {
				return err
			}

			user, err := userRepo.FindByID(req.Context(), userID)
			if err != nil {
				return apperror.Unauthorized("認証が必要です。")
			}

			ctx.Set("currentUser", user)
			return next(ctx)
		}
	}
}
