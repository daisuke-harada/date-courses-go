package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type CourseRepository interface {
	Create(ctx context.Context, course *model.Course) error
	GetByID(ctx context.Context, id uint) (*model.Course, error)
	Update(ctx context.Context, course *model.Course) error
	Delete(ctx context.Context, id uint) error
}
