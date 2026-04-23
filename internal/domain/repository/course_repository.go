package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type CourseSearchParams struct {
	PrefectureID *int
}

type CourseRepository interface {
	Create(ctx context.Context, course *model.Course) error
	FindByUserID(ctx context.Context, userID uint) ([]*model.Course, error)
	Search(ctx context.Context, params CourseSearchParams) ([]*model.Course, error)
	FindByID(ctx context.Context, id uint) (*model.Course, error)
	DeleteByID(ctx context.Context, id uint) error
}
