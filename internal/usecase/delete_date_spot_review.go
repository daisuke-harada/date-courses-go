package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type DeleteDateSpotReviewInputPort interface {
	Execute(context.Context, DeleteDateSpotReviewInput) error
}

type DeleteDateSpotReviewInput struct {
	ReviewID uint
}

type DeleteDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewDeleteDateSpotReviewUsecase(dateSpotReviewRepository repository.DateSpotReviewRepository) DeleteDateSpotReviewInputPort {
	return &DeleteDateSpotReviewInteractor{DateSpotReviewRepository: dateSpotReviewRepository}
}

func (i *DeleteDateSpotReviewInteractor) Execute(ctx context.Context, input DeleteDateSpotReviewInput) error {
	if err := i.DateSpotReviewRepository.DeleteByID(ctx, input.ReviewID); err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
