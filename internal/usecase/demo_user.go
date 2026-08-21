package usecase

import (
	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

// DemoUserName は誰でもログインできるデモ用アカウントの名前です。
// 設定から DI 経由で注入します。
type DemoUserName string

// verifyNotDemoUser はデモ用アカウントに対する変更を拒否します。
//
// デモ用アカウントは誰でもログインできる共有アカウントのため、
// プロフィールを更新されるとパスワードやメールアドレスを書き換えられ、
// 他の閲覧者がログインできなくなる。退会させられた場合も同様に復旧できない。
func verifyNotDemoUser(user *model.User, demoUserName DemoUserName) error {
	if demoUserName == "" || user.Name != string(demoUserName) {
		return nil
	}
	return apperror.Forbidden("デモ用アカウントは変更・削除できません")
}
