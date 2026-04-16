package usecase

import (
	"context"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// CreateDateSpotInputPort はデートスポット作成ユースケースの入力ポートです。
type CreateDateSpotInputPort interface {
	Execute(context.Context, CreateDateSpotInput) (*CreateDateSpotOutput, error)
}

// CreateDateSpotInput はデートスポット作成の入力データです。
type CreateDateSpotInput struct {
	Name         string
	GenreID      int
	PrefectureID int
	CityName     string
	OpeningTime  *time.Time
	ClosingTime  *time.Time
	Image        *string
}

// CreateDateSpotOutput はデートスポット作成の出力データです。
type CreateDateSpotOutput struct {
	DateSpotID uint
}

type CreateDateSpotInteractor struct {
	DateSpotRepository repository.DateSpotRepository
}

func NewCreateDateSpotUsecase(
	dateSpotRepository repository.DateSpotRepository,
) CreateDateSpotInputPort {
	return &CreateDateSpotInteractor{
		DateSpotRepository: dateSpotRepository,
	}
}

func (i *CreateDateSpotInteractor) Execute(ctx context.Context, input CreateDateSpotInput) (*CreateDateSpotOutput, error) {
	dateSpot := &model.DateSpot{
		Name:         input.Name,
		GenreID:      &input.GenreID,
		PrefectureID: &input.PrefectureID,
		CityName:     input.CityName,
		OpeningTime:  input.OpeningTime,
		ClosingTime:  input.ClosingTime,
		Image:        input.Image,
	}

	if err := i.DateSpotRepository.Create(ctx, dateSpot); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &CreateDateSpotOutput{DateSpotID: dateSpot.ID}, nil
}
