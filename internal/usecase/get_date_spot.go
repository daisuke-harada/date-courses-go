package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type GetDateSpotInputPort interface {
	Execute(context.Context, GetDateSpotInput) (*GetDateSpotOutput, error)
}

type GetDateSpotInput struct {
	ID uint
}

type GetDateSpotOutput struct {
	DateSpot *model.DateSpot
}

type GetDateSpotInteractor struct {
	DateSpotRepository repository.DateSpotRepository
}

func NewGetDateSpotUsecase(dateSpotRepository repository.DateSpotRepository) GetDateSpotInputPort {
	return &GetDateSpotInteractor{
		DateSpotRepository: dateSpotRepository,
	}
}

func (i *GetDateSpotInteractor) Execute(ctx context.Context, input GetDateSpotInput) (*GetDateSpotOutput, error) {
	dateSpot, err := i.DateSpotRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &GetDateSpotOutput{DateSpot: dateSpot}, nil
}
