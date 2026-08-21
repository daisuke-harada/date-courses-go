package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type DateSpotReviewRepository interface {
	Create(ctx context.Context, review *model.DateSpotReview) error
	FindByID(ctx context.Context, id uint) (*model.DateSpotReview, error)
	// FindByUserIDs は指定ユーザーたちのレビューを userID ごとにまとめて返します。
	FindByUserIDs(ctx context.Context, userIDs []uint) (map[uint][]*model.DateSpotReview, error)
	FindByDateSpotID(ctx context.Context, dateSpotID uint) ([]*model.DateSpotReview, error)
	DeleteByID(ctx context.Context, id uint) error
	UpdateByID(ctx context.Context, id uint, review *model.DateSpotReview) error
}
