package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repomock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteDateSpotReviewInteractor_Execute(t *testing.T) {
	// 他人のレビューを削除できないことを保証する
	t.Run("error_forbidden_when_operator_is_not_the_author", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(1)).
			Return(&model.DateSpotReview{ID: 1, UserID: 1, DateSpotID: 3}, nil)

		interactor := usecase.NewDeleteDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.DeleteDateSpotReviewInput{ReviewID: 1, OperatorID: 99})

		assert.Nil(t, output)
		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, statusCode)
	})

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(10)).
			Return(&model.DateSpotReview{ID: 10, UserID: 1, DateSpotID: 5}, nil)
		reviewRepo.EXPECT().
			DeleteByID(ctx, uint(10)).
			Return(nil)
		reviewRepo.EXPECT().
			FindByDateSpotID(ctx, uint(5)).
			Return([]*model.DateSpotReview{}, nil)

		interactor := usecase.NewDeleteDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.DeleteDateSpotReviewInput{ReviewID: 10, OperatorID: 1})

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotNil(t, output.DateSpotReviews)
	})

	t.Run("error_review_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(10)).
			Return(nil, errors.New("not found"))

		interactor := usecase.NewDeleteDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.DeleteDateSpotReviewInput{ReviewID: 10, OperatorID: 1})

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_repository_delete_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(10)).
			Return(&model.DateSpotReview{ID: 10, UserID: 1, DateSpotID: 5}, nil)
		reviewRepo.EXPECT().
			DeleteByID(ctx, uint(10)).
			Return(errors.New("db error"))

		interactor := usecase.NewDeleteDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.DeleteDateSpotReviewInput{ReviewID: 10, OperatorID: 1})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
