package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type DateSpotRepository interface {
	Create(ctx context.Context, dateSpot *model.DateSpot) error
	FindAll(ctx context.Context) ([]*model.DateSpot, error)
}
