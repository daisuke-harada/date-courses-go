package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
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

func (r *dateSpotRepository) FindByID(ctx context.Context, id uint) (*model.DateSpot, error) {
	db := r.db.WithContext(ctx).
		Model(&model.DateSpot{}).
		Select(`date_spots.*,
			COALESCE(AVG(date_spot_reviews.rate), 0)  AS average_rate,
			COUNT(date_spot_reviews.id)               AS review_total_number`).
		Joins("LEFT JOIN date_spot_reviews ON date_spot_reviews.date_spot_id = date_spots.id").
		Group("date_spots.id")

	var dateSpot model.DateSpot
	if err := db.Where("date_spots.id = ?", id).First(&dateSpot).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.NotFound()
		}
		slog.ErrorContext(ctx, "dateSpotRepository.FindByID failed", "err", err, "id", id)
		return nil, apperror.InternalServerError(err)
	}
	return &dateSpot, nil
}

func (r *dateSpotRepository) Search(ctx context.Context, params repository.DateSpotSearchParams) ([]*model.DateSpot, error) {
	db := r.db.WithContext(ctx).
		Model(&model.DateSpot{}).
		Select(`date_spots.*,
			COALESCE(AVG(date_spot_reviews.rate), 0)  AS average_rate,
			COUNT(date_spot_reviews.id)               AS review_total_number`).
		Joins("LEFT JOIN date_spot_reviews ON date_spot_reviews.date_spot_id = date_spots.id").
		Group("date_spots.id")

	if params.Name != nil && *params.Name != "" {
		db = db.Where("date_spots.name LIKE ?", "%"+*params.Name+"%")
	}
	if params.PrefectureID != nil {
		db = db.Where("date_spots.prefecture_id = ?", *params.PrefectureID)
	}
	if params.GenreID != nil {
		db = db.Where("date_spots.genre_id = ?", *params.GenreID)
	}
	if params.ComeTime != nil && *params.ComeTime != "" {
		db = db.Where("date_spots.opening_time <= ?", *params.ComeTime).
			Where("date_spots.closing_time >= ?", *params.ComeTime)
	}

	var dateSpots []*model.DateSpot
	if err := db.Find(&dateSpots).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.Search failed", "err", err)
		return nil, err
	}
	return dateSpots, nil
}

func (r *dateSpotRepository) Update(ctx context.Context, id uint, dateSpot *model.DateSpot) error {
	// Ensure the ID is set on the struct so GORM treats this as an update
	dateSpot.ID = id
	if err := r.db.WithContext(ctx).Model(&model.DateSpot{}).Where("id = ?", id).Updates(dateSpot).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.Update failed", "err", err, "id", id)
		return err
	}
	slog.InfoContext(ctx, "dateSpotRepository.Update succeeded", "id", id)
	return nil
}

func (r *dateSpotRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.DateSpot{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "dateSpotRepository.Delete failed", "err", err, "id", id)
		return err
	}
	slog.InfoContext(ctx, "dateSpotRepository.Delete succeeded", "id", id)
	return nil
}
