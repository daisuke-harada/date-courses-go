package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type CreateDateSpotReviewInputPort interface {
	Execute(context.Context, CreateDateSpotReviewInput) (*CreateDateSpotReviewOutput, error)
}

type CreateDateSpotReviewInput struct {
	UserID     uint
	DateSpotID uint
	Rate       *float64
	Content    *string
}

type CreateDateSpotReviewOutput struct {
	ReviewID uint
}

type CreateDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewCreateDateSpotReviewUsecase(dateSpotReviewRepository repository.DateSpotReviewRepository) CreateDateSpotReviewInputPort {
	return &CreateDateSpotReviewInteractor{DateSpotReviewRepository: dateSpotReviewRepository}
}

func (i *CreateDateSpotReviewInteractor) Execute(ctx context.Context, input CreateDateSpotReviewInput) (*CreateDateSpotReviewOutput, error) {
	review := &model.DateSpotReview{
		UserID:     input.UserID,
		DateSpotID: input.DateSpotID,
		Rate:       input.Rate,
		Content:    input.Content,
	}
	if err := i.DateSpotReviewRepository.Create(ctx, review); err != nil {
		return nil, apperror.InternalServerError(err)
	}
	return &CreateDateSpotReviewOutput{ReviewID: review.ID}, nil
}
