package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type DateSpotReviewRepository interface {
	Create(ctx context.Context, review *model.DateSpotReview) error
	FindByUserID(ctx context.Context, userID uint) ([]*model.DateSpotReview, error)
	FindByID(ctx context.Context, id uint) (*model.DateSpotReview, error)
	Update(ctx context.Context, review *model.DateSpotReview) error
	DeleteByID(ctx context.Context, id uint) error
}
