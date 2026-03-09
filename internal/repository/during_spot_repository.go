package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
)

type DuringSpotRepository interface {
	Create(ctx context.Context, duringSpot *domain.DuringSpot) error
	GetByID(ctx context.Context, id uint) (*domain.DuringSpot, error)
	Update(ctx context.Context, duringSpot *domain.DuringSpot) error
	Delete(ctx context.Context, id uint) error
}
