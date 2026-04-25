package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

// DateSpotSearchParams はdate_spotsの検索条件を表します。
type DateSpotSearchParams struct {
	Name         *string
	PrefectureID *int
	GenreID      *int
	ComeTime     *string
}

type DateSpotRepository interface {
	Create(ctx context.Context, dateSpot *model.DateSpot) error
	FindByID(ctx context.Context, id uint) (*model.DateSpot, error)
	Search(ctx context.Context, params DateSpotSearchParams) ([]*model.DateSpot, error)
	Update(ctx context.Context, id uint, dateSpot *model.DateSpot) error
	Delete(ctx context.Context, id uint) error
}
