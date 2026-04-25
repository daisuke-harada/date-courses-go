package middleware

import (
	"net/http"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	jwtpkg "github.com/daisuke-harada/date-courses-go/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
)

// publicPaths は認証不要なパスのセットです。
// POST /api/v1/login, POST /api/v1/signup はメソッドも一致させて除外します。
// GET 系の公開エンドポイントはメソッドがGETの場合のみ除外します。
var publicGETPrefixes = []string{
	"/",
	"/api/v1/top",
	"/api/v1/date_spots",
	"/api/v1/date_spot_reviews",
	"/api/v1/prefectures",
	"/api/v1/genres",
	"/api/v1/users",
	"/api/v1/courses",
}

var publicPOSTPaths = []string{
	"/api/v1/login",
	"/api/v1/signup",
}

// protectedGETSuffixes は GET であっても認証が必要なパスのサフィックスです。
var protectedGETSuffixes = []string{
	"/followings",
	"/followers",
}

func isPublic(method, path string) bool {
	if method == http.MethodPost {
		for _, p := range publicPOSTPaths {
			if path == p {
				return true
			}
		}
		return false
	}
	if method == http.MethodGet {
		for _, suffix := range protectedGETSuffixes {
			if strings.HasSuffix(path, suffix) {
				return false
			}
		}
		for _, prefix := range publicGETPrefixes {
			if path == prefix || strings.HasPrefix(path, prefix+"/") || strings.HasPrefix(path, prefix+"?") {
				return true
			}
		}
	}
	return false
}

// JWTAuthMiddleware は JWT Bearer トークンを検証し、currentUser をコンテキストにセットします。
func JWTAuthMiddleware(secretKey string, userRepo repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			if isPublic(req.Method, req.URL.Path) {
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
