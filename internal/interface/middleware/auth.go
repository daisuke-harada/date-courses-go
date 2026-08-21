package middleware

import (
	"net/http"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	iface_openapi "github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	jwtpkg "github.com/daisuke-harada/date-courses-go/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
)

// JWTAuthMiddleware は JWT Bearer トークンを検証し、currentUser をコンテキストにセットします。
// 認証が必要かどうかは iface_openapi.RequiresBearerAuth を通じて判定します。
// この関数は api/resolved/openapi/openapi.yaml の security 定義から自動生成された
// auth_routes.gen.go のマップを参照します。`make gen` で再生成されます。
//
// 認証必須ではないルートでも、トークンが付いていれば currentUser を解決します（任意認証）。
// 非公開のデートコースを作成者にだけ返すために、公開 GET でも「誰が見ているか」が必要なためです。
// 任意認証のルートではトークンが無効でも 401 にせず匿名として扱います。401 にすると、
// 期限切れトークンを持ったままログイン画面を開いたユーザーがログイン自体を拒否され、
// 再ログインできなくなるためです。
func JWTAuthMiddleware(secretKey string, userRepo repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			user, err := authenticate(req, secretKey, userRepo)
			if err != nil {
				if iface_openapi.RequiresBearerAuth(req.Method, ctx.Path()) {
					return err
				}
				return next(ctx)
			}

			SetCurrentUser(ctx, user)
			return next(ctx)
		}
	}
}

// authenticate は Authorization ヘッダの Bearer トークンを検証し、対応するユーザーを返します。
func authenticate(req *http.Request, secretKey string, userRepo repository.UserRepository) (*model.User, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, apperror.Unauthorized("認証が必要です。")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := jwtpkg.Decode(tokenStr, secretKey)
	if err != nil {
		return nil, err
	}

	user, err := userRepo.FindByID(req.Context(), userID)
	if err != nil {
		return nil, apperror.Unauthorized("認証が必要です。")
	}

	return user, nil
}
