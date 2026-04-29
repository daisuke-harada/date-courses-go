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
	DateSpot        *model.DateSpot
	DateSpotReviews []*model.DateSpotReview
}

type GetDateSpotInteractor struct {
	DateSpotRepository       repository.DateSpotRepository
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewGetDateSpotUsecase(
	dateSpotRepository repository.DateSpotRepository,
	dateSpotReviewRepository repository.DateSpotReviewRepository,
) GetDateSpotInputPort {
	return &GetDateSpotInteractor{
		DateSpotRepository:       dateSpotRepository,
		DateSpotReviewRepository: dateSpotReviewRepository,
	}
}

func (i *GetDateSpotInteractor) Execute(ctx context.Context, input GetDateSpotInput) (*GetDateSpotOutput, error) {
	dateSpot, err := i.DateSpotRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	reviews, err := i.DateSpotReviewRepository.FindByDateSpotID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetDateSpotOutput{DateSpot: dateSpot, DateSpotReviews: reviews}, nil
}
