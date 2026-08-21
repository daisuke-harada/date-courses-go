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

// CurrentUserID は認証済みユーザーの ID を返します。未認証の場合は 0 を返します。
// 非公開データの絞り込みなど「誰が見ているか」だけが必要な場面で使います。
func CurrentUserID(ctx echo.Context) uint {
	user := CurrentUser(ctx)
	if user == nil {
		return 0
	}
	return user.ID
}

// RequireAdmin は管理者を返します。未認証なら 401、管理者でなければ 403 を返します。
// デートスポットの登録・編集・削除など、管理者だけに許す操作で使います。
func RequireAdmin(ctx echo.Context) (*model.User, error) {
	user, err := RequireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if !user.Admin {
		return nil, apperror.Forbidden("管理者のみ実行できます")
	}
	return user, nil
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
