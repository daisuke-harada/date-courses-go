package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type DuringSpotRepository interface {
	Create(ctx context.Context, duringSpot *model.DuringSpot) error
	GetByID(ctx context.Context, id uint) (*model.DuringSpot, error)
	Update(ctx context.Context, duringSpot *model.DuringSpot) error
	Delete(ctx context.Context, id uint) error
}
