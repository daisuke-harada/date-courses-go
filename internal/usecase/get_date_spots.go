package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
)

type GetDateSpotsInputPort interface {
	Execute(context.Context) (GetDateSpotsOutput, error)
}

type GetDateSpotsOutput struct {
	DateSpots []*model.DateSpot
}

type GetDateSpotsInteractor struct {
	DateSpotRepository repository.DateSpotRepository
}

func NewGetDateSpotsUsecase(
	dateSpotRepository repository.DateSpotRepository,
) GetDateSpotsInputPort {
	return &GetDateSpotsInteractor{
		DateSpotRepository: dateSpotRepository,
	}
}

func (i *GetDateSpotsInteractor) Execute(ctx context.Context) (GetDateSpotsOutput, error) {
	dateSpots, err := i.DateSpotRepository.FindAll(ctx)
	if err != nil {
		return GetDateSpotsOutput{}, err
	}

	return GetDateSpotsOutput{
		DateSpots: dateSpots,
	}, nil
}
