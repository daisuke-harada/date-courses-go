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
	// OperatorID は削除を実行するユーザー（トークンの currentUser）の ID です。
	OperatorID uint
}

type DeleteDateSpotReviewOutput struct {
	DateSpotReviews []*model.DateSpotReview
}

type DeleteDateSpotReviewInteractor struct {
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewDeleteDateSpotReviewUsecase(dateSpotReviewRepository repository.DateSpotReviewRepository) DeleteDateSpotReviewInputPort {
	return &DeleteDateSpotReviewInteractor{DateSpotReviewRepository: dateSpotReviewRepository}
}

func (i *DeleteDateSpotReviewInteractor) Execute(ctx context.Context, input DeleteDateSpotReviewInput) (*DeleteDateSpotReviewOutput, error) {
	review, err := i.DateSpotReviewRepository.FindByID(ctx, input.ReviewID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	// レビューを削除できるのは投稿者だけ
	if review.UserID != input.OperatorID {
		return nil, apperror.Forbidden("他のユーザーのレビューは削除できません")
	}

	dateSpotID := review.DateSpotID

	if err := i.DateSpotReviewRepository.DeleteByID(ctx, input.ReviewID); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	reviews, err := i.DateSpotReviewRepository.FindByDateSpotID(ctx, dateSpotID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &DeleteDateSpotReviewOutput{
		DateSpotReviews: reviews,
	}, nil
}
