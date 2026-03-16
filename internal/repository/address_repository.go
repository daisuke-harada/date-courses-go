package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type AddressRepository interface {
	Create(ctx context.Context, address *model.Address) error
	GetByID(ctx context.Context, id uint) (*model.Address, error)
	Update(ctx context.Context, address *model.Address) error
	Delete(ctx context.Context, id uint) error
	FindByDateSpotID(ctx context.Context, dateSpotID uint) (*model.Address, error)
}
