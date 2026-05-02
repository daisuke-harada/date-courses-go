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

func (i *CreateDateSpotReviewInput) Validate() error {
	var errs []string
	if i.UserID == 0 {
		errs = append(errs, "ユーザーIDを入力してください")
	}
	if i.DateSpotID == 0 {
		errs = append(errs, "デートスポットIDを入力してください")
	}
	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}
	return nil
}

type CreateDateSpotReviewOutput struct {
	ReviewID        uint
	DateSpotReviews []*model.DateSpotReview
}

type CreateDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewCreateDateSpotReviewUsecase(dateSpotReviewRepository repository.DateSpotReviewRepository) CreateDateSpotReviewInputPort {
	return &CreateDateSpotReviewInteractor{DateSpotReviewRepository: dateSpotReviewRepository}
}

func (i *CreateDateSpotReviewInteractor) Execute(ctx context.Context, input CreateDateSpotReviewInput) (*CreateDateSpotReviewOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	review := &model.DateSpotReview{
		UserID:     input.UserID,
		DateSpotID: input.DateSpotID,
		Rate:       input.Rate,
		Content:    input.Content,
	}
	if err := i.DateSpotReviewRepository.Create(ctx, review); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	reviews, err := i.DateSpotReviewRepository.FindByDateSpotID(ctx, input.DateSpotID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &CreateDateSpotReviewOutput{
		ReviewID:        review.ID,
		DateSpotReviews: reviews,
	}, nil
}
