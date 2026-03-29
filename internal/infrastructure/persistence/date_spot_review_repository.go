package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
	"gorm.io/gorm"
)

type dateSpotReviewRepository struct {
	db *gorm.DB
}

func NewDateSpotReviewRepository(db *gorm.DB) repository.DateSpotReviewRepository {
	return &dateSpotReviewRepository{db: db}
}

func (r *dateSpotReviewRepository) Create(ctx context.Context, review *model.DateSpotReview) error {
	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "dateSpotReviewRepository.Create succeeded", "review_id", review.ID)
	return nil
}
