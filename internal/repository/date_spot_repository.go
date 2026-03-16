package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type DateSpotRepository interface {
	Create(ctx context.Context, dateSpot *model.DateSpot) error
	GetByID(ctx context.Context, id uint) (*model.DateSpot, error)
	Update(ctx context.Context, dateSpot *model.DateSpot) error
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]*model.DateSpot, error)
}
