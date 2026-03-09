package repository

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
	"gorm.io/gorm"
)

type DuringSpotRepository interface {
	Create(ctx context.Context, duringSpot *domain.DuringSpot) error
	GetByID(ctx context.Context, id uint) (*domain.DuringSpot, error)
	Update(ctx context.Context, duringSpot *domain.DuringSpot) error
	Delete(ctx context.Context, id uint) error
}

type duringSpotRepository struct {
	db *gorm.DB
}

func NewDuringSpotRepository(db *gorm.DB) DuringSpotRepository {
	return &duringSpotRepository{db: db}
}

func (r *duringSpotRepository) Create(ctx context.Context, duringSpot *domain.DuringSpot) error {
	if err := r.db.WithContext(ctx).Create(duringSpot).Error; err != nil {
		slog.ErrorContext(ctx, "duringSpotRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "duringSpotRepository.Create succeeded", "during_spot_id", duringSpot.ID)
	return nil
}

func (r *duringSpotRepository) GetByID(ctx context.Context, id uint) (*domain.DuringSpot, error) {
	var duringSpot domain.DuringSpot
	if err := r.db.WithContext(ctx).First(&duringSpot, id).Error; err != nil {
		slog.ErrorContext(ctx, "duringSpotRepository.GetByID failed", "during_spot_id", id, "err", err)
		return nil, err
	}
	return &duringSpot, nil
}

func (r *duringSpotRepository) Update(ctx context.Context, duringSpot *domain.DuringSpot) error {
	if err := r.db.WithContext(ctx).Save(duringSpot).Error; err != nil {
		slog.ErrorContext(ctx, "duringSpotRepository.Update failed", "during_spot_id", duringSpot.ID, "err", err)
		return err
	}
	slog.InfoContext(ctx, "duringSpotRepository.Update succeeded", "during_spot_id", duringSpot.ID)
	return nil
}

func (r *duringSpotRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.DuringSpot{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "duringSpotRepository.Delete failed", "during_spot_id", id, "err", err)
		return err
	}
	slog.InfoContext(ctx, "duringSpotRepository.Delete succeeded", "during_spot_id", id)
	return nil
}
