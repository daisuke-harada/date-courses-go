package usecase

import (
	"context"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// UpdateDateSpotInputPort はデートスポット更新ユースケースの入力ポートです。
type UpdateDateSpotInputPort interface {
	Execute(context.Context, UpdateDateSpotInput) error
}

// UpdateDateSpotInput はデートスポット更新の入力データです。
type UpdateDateSpotInput struct {
	DateSpotID   uint
	Name         string
	GenreID      int
	PrefectureID int
	CityName     string
	OpeningTime  *time.Time
	ClosingTime  *time.Time
	Image        *string
}

type UpdateDateSpotInteractor struct {
	DateSpotRepository repository.DateSpotRepository
}

func NewUpdateDateSpotUsecase(
	dateSpotRepository repository.DateSpotRepository,
) UpdateDateSpotInputPort {
	return &UpdateDateSpotInteractor{
		DateSpotRepository: dateSpotRepository,
	}
}

func (i *UpdateDateSpotInteractor) Execute(ctx context.Context, input UpdateDateSpotInput) error {
	dateSpot := &model.DateSpot{
		Name:         input.Name,
		GenreID:      &input.GenreID,
		PrefectureID: &input.PrefectureID,
		CityName:     input.CityName,
		OpeningTime:  input.OpeningTime,
		ClosingTime:  input.ClosingTime,
		Image:        input.Image,
	}

	if err := i.DateSpotRepository.Update(ctx, input.DateSpotID, dateSpot); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}
