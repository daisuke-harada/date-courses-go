package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type DuringSpotRepository interface {
	Create(ctx context.Context, duringSpot *model.DuringSpot) error
}
