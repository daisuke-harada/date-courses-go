package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
	"gorm.io/gorm"
)

type dateSpotReviewRepository struct {
	db *gorm.DB
}

func NewDateSpotReviewRepository(db *gorm.DB) repository.DateSpotReviewRepository {
	return &dateSpotReviewRepository{db: db}
}

func (r *dateSpotReviewRepository) Create(ctx context.Context, review *domain.DateSpotReview) error {
	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "dateSpotReviewRepository.Create succeeded", "review_id", review.ID)
	return nil
}

func (r *dateSpotReviewRepository) GetByID(ctx context.Context, id uint) (*domain.DateSpotReview, error) {
	var review domain.DateSpotReview
	if err := r.db.WithContext(ctx).First(&review, id).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.GetByID failed", "review_id", id, "err", err)
		return nil, err
	}
	return &review, nil
}

func (r *dateSpotReviewRepository) Update(ctx context.Context, review *domain.DateSpotReview) error {
	if err := r.db.WithContext(ctx).Save(review).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.Update failed", "review_id", review.ID, "err", err)
		return err
	}
	slog.InfoContext(ctx, "dateSpotReviewRepository.Update succeeded", "review_id", review.ID)
	return nil
}

func (r *dateSpotReviewRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.DateSpotReview{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.Delete failed", "review_id", id, "err", err)
		return err
	}
	slog.InfoContext(ctx, "dateSpotReviewRepository.Delete succeeded", "review_id", id)
	return nil
}
