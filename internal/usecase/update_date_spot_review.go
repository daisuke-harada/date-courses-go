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
	ReviewID uint
	Rate     *float64
	Content  *string
}

// Validate はレビュー更新の入力バリデーションを行います。
// rate または content のいずれかが必要です。
func (i *UpdateDateSpotReviewInput) Validate() error {
	if i.Rate == nil && i.Content == nil {
		return apperror.UnprocessableEntity("rate または content のいずれかを入力してください")
	}
	return nil
}

type UpdateDateSpotReviewOutput struct {
	ReviewID        uint
	DateSpot        *model.DateSpot
	DateSpotReviews []*model.DateSpotReview
}

type UpdateDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
	DateSpotRepository       repository.DateSpotRepository
}

func NewUpdateDateSpotReviewUsecase(
	dateSpotReviewRepository repository.DateSpotReviewRepository,
	dateSpotRepository repository.DateSpotRepository,
) UpdateDateSpotReviewInputPort {
	return &UpdateDateSpotReviewInteractor{
		DateSpotReviewRepository: dateSpotReviewRepository,
		DateSpotRepository:       dateSpotRepository,
	}
}

func (i *UpdateDateSpotReviewInteractor) Execute(ctx context.Context, input UpdateDateSpotReviewInput) (*UpdateDateSpotReviewOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	existing, err := i.DateSpotReviewRepository.FindByID(ctx, input.ReviewID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	review := &model.DateSpotReview{
		Rate:    input.Rate,
		Content: input.Content,
	}
	if err := i.DateSpotReviewRepository.UpdateByID(ctx, input.ReviewID, review); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	dateSpot, err := i.DateSpotRepository.FindByID(ctx, existing.DateSpotID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	reviews, err := i.DateSpotReviewRepository.FindByDateSpotID(ctx, existing.DateSpotID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &UpdateDateSpotReviewOutput{
		ReviewID:        input.ReviewID,
		DateSpot:        dateSpot,
		DateSpotReviews: reviews,
	}, nil
}
