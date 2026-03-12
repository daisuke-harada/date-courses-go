package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
)

type DateSpotReviewRepository interface {
	Create(ctx context.Context, review *domain.DateSpotReview) error
	GetByID(ctx context.Context, id uint) (*domain.DateSpotReview, error)
	Update(ctx context.Context, review *domain.DateSpotReview) error
	Delete(ctx context.Context, id uint) error
}
