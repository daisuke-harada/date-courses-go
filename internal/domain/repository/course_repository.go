package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type CourseRepository interface {
	Create(ctx context.Context, course *model.Course) error
	FindByUserID(ctx context.Context, userID uint) ([]*model.Course, error)
	FindAll(ctx context.Context) ([]*model.Course, error)
	FindByID(ctx context.Context, id uint) (*model.Course, error)
}
