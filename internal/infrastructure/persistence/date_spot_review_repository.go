package persistence

import (
	"context"
	"log/slog"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
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

// DeleteByID は指定 ID のレビューを削除します。
func (r *dateSpotReviewRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.DateSpotReview{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.DeleteByID failed", "err", err)
		return err
	}
	return nil
}

// FindByUserID は指定ユーザーのレビュー一覧を DateSpot 込みで返します。
func (r *dateSpotReviewRepository) FindByUserID(ctx context.Context, userID uint) ([]*model.DateSpotReview, error) {
	var reviews []*model.DateSpotReview
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("DateSpot").
		Find(&reviews).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.FindByUserID failed", "err", err)
		return nil, err
	}
	return reviews, nil
}
