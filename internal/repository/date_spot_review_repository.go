package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type DateSpotReviewRepository interface {
	Create(ctx context.Context, review *model.DateSpotReview) error
	GetByID(ctx context.Context, id uint) (*model.DateSpotReview, error)
	Update(ctx context.Context, review *model.DateSpotReview) error
	Delete(ctx context.Context, id uint) error
}
