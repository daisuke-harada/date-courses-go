package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type UpdateDateSpotReviewInputPort interface {
	Execute(context.Context, UpdateDateSpotReviewInput) (*UpdateDateSpotReviewOutput, error)
}

type UpdateDateSpotReviewInput struct {
	ReviewID uint
	Rate     *float64
	Content  *string
}

type UpdateDateSpotReviewOutput struct{}

type UpdateDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewUpdateDateSpotReviewUsecase(dateSpotReviewRepository repository.DateSpotReviewRepository) UpdateDateSpotReviewInputPort {
	return &UpdateDateSpotReviewInteractor{DateSpotReviewRepository: dateSpotReviewRepository}
}

func (i *UpdateDateSpotReviewInteractor) Execute(ctx context.Context, input UpdateDateSpotReviewInput) (*UpdateDateSpotReviewOutput, error) {
	review, err := i.DateSpotReviewRepository.FindByID(ctx, input.ReviewID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	review.Rate = input.Rate
	review.Content = input.Content

	if err := i.DateSpotReviewRepository.Update(ctx, review); err != nil {
		return nil, apperror.InternalServerError(err)
	}
	return &UpdateDateSpotReviewOutput{}, nil
}
