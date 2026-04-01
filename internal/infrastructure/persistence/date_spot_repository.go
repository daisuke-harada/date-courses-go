package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"gorm.io/gorm"
)

type dateSpotRepository struct {
	db *gorm.DB
}

func NewDateSpotRepository(db *gorm.DB) repository.DateSpotRepository {
	return &dateSpotRepository{db: db}
}

func (r *dateSpotRepository) Create(ctx context.Context, dateSpot *model.DateSpot) error {
	if err := r.db.WithContext(ctx).Create(dateSpot).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "dateSpotRepository.Create succeeded", "date_spot_id", dateSpot.ID)
	return nil
}

func (r *dateSpotRepository) Search(ctx context.Context, params repository.DateSpotSearchParams) ([]*model.DateSpot, error) {
	db := r.db.WithContext(ctx).Model(&model.DateSpot{})

	if params.Name != nil && *params.Name != "" {
		db = db.Where("name LIKE ?", "%"+*params.Name+"%")
	}
	if params.PrefectureID != nil {
		db = db.Where("prefecture_id = ?", *params.PrefectureID)
	}
	if params.GenreID != nil {
		db = db.Where("genre_id = ?", *params.GenreID)
	}
	if params.ComeTime != nil && *params.ComeTime != "" {
		db = db.Where("opening_time <= ?", *params.ComeTime).
			Where("closing_time >= ?", *params.ComeTime)
	}

	var dateSpots []*model.DateSpot
	if err := db.Find(&dateSpots).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.Search failed", "err", err)
		return nil, err
	}
	return dateSpots, nil
}
