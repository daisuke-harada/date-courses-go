package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type CourseRepository interface {
	Create(ctx context.Context, course *model.Course) error
}
