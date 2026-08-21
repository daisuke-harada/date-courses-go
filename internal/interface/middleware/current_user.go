package middleware

import (
	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/labstack/echo/v4"
)

// currentUserKey は認証済みユーザーを echo.Context に格納する際のキーです。
// 格納と取り出しを同じパッケージに閉じるため非公開にしています。
const currentUserKey = "currentUser"

// SetCurrentUser は認証済みユーザーを echo.Context に格納します。
func SetCurrentUser(ctx echo.Context, user *model.User) {
	ctx.Set(currentUserKey, user)
}

// CurrentUser は JWTAuthMiddleware がセットした認証済みユーザーを返します。
// 未認証（任意認証ルートにトークンなしでアクセスした場合など）は nil を返します。
func CurrentUser(ctx echo.Context) *model.User {
	user, ok := ctx.Get(currentUserKey).(*model.User)
	if !ok {
		return nil
	}
	return user
}

// RequireCurrentUser は認証済みユーザーを返します。未認証の場合は 401 を返します。
// 認証必須ルートでも、リクエストの user_id ではなく必ずこのユーザーを操作主体として扱います。
func RequireCurrentUser(ctx echo.Context) (*model.User, error) {
	user := CurrentUser(ctx)
	if user == nil {
		return nil, apperror.Unauthorized("認証が必要です。")
	}
	return user, nil
}
