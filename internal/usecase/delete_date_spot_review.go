package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type DeleteDateSpotReviewInputPort interface {
	Execute(context.Context, DeleteDateSpotReviewInput) (*DeleteDateSpotReviewOutput, error)
}

type DeleteDateSpotReviewInput struct {
	ReviewID uint
}

type DeleteDateSpotReviewOutput struct {
	DateSpot        *model.DateSpot
	DateSpotReviews []*model.DateSpotReview
}

type DeleteDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
	DateSpotRepository       repository.DateSpotRepository
}

func NewDeleteDateSpotReviewUsecase(
	dateSpotReviewRepository repository.DateSpotReviewRepository,
	dateSpotRepository repository.DateSpotRepository,
) DeleteDateSpotReviewInputPort {
	return &DeleteDateSpotReviewInteractor{
		DateSpotReviewRepository: dateSpotReviewRepository,
		DateSpotRepository:       dateSpotRepository,
	}
}

func (i *DeleteDateSpotReviewInteractor) Execute(ctx context.Context, input DeleteDateSpotReviewInput) (*DeleteDateSpotReviewOutput, error) {
	review, err := i.DateSpotReviewRepository.FindByID(ctx, input.ReviewID)
	if err != nil {
		return nil, apperror.NotFound()
	}
	dateSpotID := review.DateSpotID

	if err := i.DateSpotReviewRepository.DeleteByID(ctx, input.ReviewID); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	dateSpot, err := i.DateSpotRepository.FindByID(ctx, dateSpotID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	reviews, err := i.DateSpotReviewRepository.FindByDateSpotID(ctx, dateSpotID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &DeleteDateSpotReviewOutput{
		DateSpot:        dateSpot,
		DateSpotReviews: reviews,
	}, nil
}
