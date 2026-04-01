package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type DateSpotReviewRepository interface {
	Create(ctx context.Context, review *model.DateSpotReview) error
}
