package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
)

type CourseRepository interface {
	Create(ctx context.Context, course *domain.Course) error
	GetByID(ctx context.Context, id uint) (*domain.Course, error)
	Update(ctx context.Context, course *domain.Course) error
	Delete(ctx context.Context, id uint) error
}
