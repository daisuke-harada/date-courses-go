package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
	"gorm.io/gorm"
)

type dateSpotRepository struct {
	db *gorm.DB
}

func NewDateSpotRepository(db *gorm.DB) repository.DateSpotRepository {
	return &dateSpotRepository{db: db}
}

func (r *dateSpotRepository) Create(ctx context.Context, dateSpot *domain.DateSpot) error {
	if err := r.db.WithContext(ctx).Create(dateSpot).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "dateSpotRepository.Create succeeded", "date_spot_id", dateSpot.ID)
	return nil
}

func (r *dateSpotRepository) GetByID(ctx context.Context, id uint) (*domain.DateSpot, error) {
	var dateSpot domain.DateSpot
	if err := r.db.WithContext(ctx).First(&dateSpot, id).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.GetByID failed", "date_spot_id", id, "err", err)
		return nil, err
	}
	return &dateSpot, nil
}

func (r *dateSpotRepository) Update(ctx context.Context, dateSpot *domain.DateSpot) error {
	if err := r.db.WithContext(ctx).Save(dateSpot).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.Update failed", "date_spot_id", dateSpot.ID, "err", err)
		return err
	}
	slog.InfoContext(ctx, "dateSpotRepository.Update succeeded", "date_spot_id", dateSpot.ID)
	return nil
}

func (r *dateSpotRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.DateSpot{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.Delete failed", "date_spot_id", id, "err", err)
		return err
	}
	slog.InfoContext(ctx, "dateSpotRepository.Delete succeeded", "date_spot_id", id)
	return nil
}
