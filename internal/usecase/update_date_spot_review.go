package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type UpdateDateSpotReviewInputPort interface {
	Execute(context.Context, UpdateDateSpotReviewInput) (*UpdateDateSpotReviewOutput, error)
}

type UpdateDateSpotReviewInput struct {
	ReviewID   uint
	DateSpotID uint
	Rate       *float64
	Content    *string
}

func (i *UpdateDateSpotReviewInput) Validate() error {
	if i.Rate == nil && i.Content == nil {
		return apperror.UnprocessableEntity("rate または content のいずれかを入力してください")
	}
	return nil
}

type UpdateDateSpotReviewOutput struct {
	ReviewID        uint
	DateSpotReviews []*model.DateSpotReview
}

type UpdateDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewUpdateDateSpotReviewUsecase(dateSpotReviewRepository repository.DateSpotReviewRepository) UpdateDateSpotReviewInputPort {
	return &UpdateDateSpotReviewInteractor{DateSpotReviewRepository: dateSpotReviewRepository}
}

func (i *UpdateDateSpotReviewInteractor) Execute(ctx context.Context, input UpdateDateSpotReviewInput) (*UpdateDateSpotReviewOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	review := &model.DateSpotReview{
		Rate:    input.Rate,
		Content: input.Content,
	}
	if err := i.DateSpotReviewRepository.UpdateByID(ctx, input.ReviewID, review); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	reviews, err := i.DateSpotReviewRepository.FindByDateSpotID(ctx, input.DateSpotID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &UpdateDateSpotReviewOutput{
		ReviewID:        input.ReviewID,
		DateSpotReviews: reviews,
	}, nil
}
