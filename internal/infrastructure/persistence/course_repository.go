package persistence

import (
	"context"
	"log/slog"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"gorm.io/gorm"
)

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) repository.CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, course *model.Course) error {
	if err := r.db.WithContext(ctx).Create(course).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "courseRepository.Create succeeded", "course_id", course.ID)
	return nil
}

// FindByUserID は指定ユーザーのコース一覧を DuringSpots→DateSpot 込みで返します。
func (r *courseRepository) FindByUserID(ctx context.Context, userID uint) ([]*model.Course, error) {
	var courses []*model.Course
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("DuringSpots.DateSpot").
		Find(&courses).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.FindByUserID failed", "err", err)
		return nil, err
	}
	return courses, nil
}

// FindAll はすべてのコース一覧を返します。
func (r *courseRepository) FindAll(ctx context.Context) ([]*model.Course, error) {
	var courses []*model.Course
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("DuringSpots.DateSpot").
		Find(&courses).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.FindAll failed", "err", err)
		return nil, err
	}
	return courses, nil
}
