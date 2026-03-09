package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
	"gorm.io/gorm"
)

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) repository.CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, course *domain.Course) error {
	if err := r.db.WithContext(ctx).Create(course).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "courseRepository.Create succeeded", "course_id", course.ID)
	return nil
}

func (r *courseRepository) GetByID(ctx context.Context, id uint) (*domain.Course, error) {
	var course domain.Course
	if err := r.db.WithContext(ctx).First(&course, id).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.GetByID failed", "course_id", id, "err", err)
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(ctx context.Context, course *domain.Course) error {
	if err := r.db.WithContext(ctx).Save(course).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.Update failed", "course_id", course.ID, "err", err)
		return err
	}
	slog.InfoContext(ctx, "courseRepository.Update succeeded", "course_id", course.ID)
	return nil
}

func (r *courseRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Course{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "courseRepository.Delete failed", "course_id", id, "err", err)
		return err
	}
	slog.InfoContext(ctx, "courseRepository.Delete succeeded", "course_id", id)
	return nil
}
