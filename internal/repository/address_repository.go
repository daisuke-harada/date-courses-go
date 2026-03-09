package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
)

type AddressRepository interface {
	Create(ctx context.Context, address *domain.Address) error
	GetByID(ctx context.Context, id uint) (*domain.Address, error)
	Update(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id uint) error
}
