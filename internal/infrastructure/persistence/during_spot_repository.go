package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"gorm.io/gorm"
)

type duringSpotRepository struct {
	db *gorm.DB
}

func NewDuringSpotRepository(db *gorm.DB) repository.DuringSpotRepository {
	return &duringSpotRepository{db: db}
}

func (r *duringSpotRepository) Create(ctx context.Context, duringSpot *model.DuringSpot) error {
	if err := r.db.WithContext(ctx).Create(duringSpot).Error; err != nil {
		slog.ErrorContext(ctx, "duringSpotRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "duringSpotRepository.Create succeeded", "during_spot_id", duringSpot.ID)
	return nil
}
