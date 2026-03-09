package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
)

type DateSpotRepository interface {
	Create(ctx context.Context, dateSpot *domain.DateSpot) error
	GetByID(ctx context.Context, id uint) (*domain.DateSpot, error)
	Update(ctx context.Context, dateSpot *domain.DateSpot) error
	Delete(ctx context.Context, id uint) error
}
