package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type CourseSearchParams struct {
	PrefectureID *int
}

// 非公開コースは作成者本人にしか見せないため、取得系はデフォルトで公開コースだけを返します。
// 非公開コースを含めるのは、本人がマイページを開いたときに使う FindAllByUserID だけです。
type CourseRepository interface {
	Create(ctx context.Context, course *model.Course) error
	// FindPublicByUserID は指定ユーザーの公開コースだけを返します。
	FindPublicByUserID(ctx context.Context, userID uint) ([]*model.Course, error)
	// FindAllByUserID は非公開コースも含めて返します。本人のマイページ専用です。
	FindAllByUserID(ctx context.Context, userID uint) ([]*model.Course, error)
	Search(ctx context.Context, params CourseSearchParams) ([]*model.Course, error)
	// FindByID は公開コース、または viewerID 自身が作成した非公開コースを返します。
	// viewerID は閲覧しているユーザーの ID で、未ログインの場合は 0 を渡します。
	FindByID(ctx context.Context, id, viewerID uint) (*model.Course, error)
	DeleteByID(ctx context.Context, id uint) error
}
