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

// FindByDateSpotID は指定 DateSpot のレビュー一覧を User 込みで返します。
func (r *dateSpotReviewRepository) FindByDateSpotID(ctx context.Context, dateSpotID uint) ([]*model.DateSpotReview, error) {
	var reviews []*model.DateSpotReview
	if err := r.db.WithContext(ctx).
		Where("date_spot_id = ?", dateSpotID).
		Preload("User").
		Find(&reviews).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.FindByDateSpotID failed", "err", err)
		return nil, err
	}
	return reviews, nil
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

// UpdateByID は指定 ID のレビューを更新します。nil フィールドは更新しません。
func (r *dateSpotReviewRepository) UpdateByID(ctx context.Context, id uint, review *model.DateSpotReview) error {
	updates := map[string]interface{}{}
	if review.Rate != nil {
		updates["rate"] = review.Rate
	}
	if review.Content != nil {
		updates["content"] = review.Content
	}
	if err := r.db.WithContext(ctx).Model(&model.DateSpotReview{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotReviewRepository.UpdateByID failed", "err", err)
		return err
	}
	return nil
}
