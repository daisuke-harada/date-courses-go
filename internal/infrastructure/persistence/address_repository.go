package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) repository.AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *model.Address) error {
	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		slog.ErrorContext(ctx, "addressRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "addressRepository.Create succeeded", "address_id", address.ID)
	return nil
}

func (r *addressRepository) GetByID(ctx context.Context, id uint) (*model.Address, error) {
	var address model.Address
	if err := r.db.WithContext(ctx).First(&address, id).Error; err != nil {
		slog.ErrorContext(ctx, "addressRepository.GetByID failed", "address_id", id, "err", err)
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) Update(ctx context.Context, address *model.Address) error {
	if err := r.db.WithContext(ctx).Save(address).Error; err != nil {
		slog.ErrorContext(ctx, "addressRepository.Update failed", "address_id", address.ID, "err", err)
		return err
	}
	slog.InfoContext(ctx, "addressRepository.Update succeeded", "address_id", address.ID)
	return nil
}

func (r *addressRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Address{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "addressRepository.Delete failed", "address_id", id, "err", err)
		return err
	}
	slog.InfoContext(ctx, "addressRepository.Delete succeeded", "address_id", id)
	return nil
}

func (r *addressRepository) FindByDateSpotID(ctx context.Context, dateSpotID uint) (*model.Address, error) {
	var address model.Address
	if err := r.db.WithContext(ctx).Where("date_spot_id = ?", dateSpotID).First(&address).Error; err != nil {
		slog.ErrorContext(ctx, "addressRepository.FindByDateSpotID failed", "date_spot_id", dateSpotID, "err", err)
		return nil, err
	}
	return &address, nil
}
