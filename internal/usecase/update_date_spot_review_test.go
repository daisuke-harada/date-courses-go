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

func TestUpdateDateSpotReviewInteractor_Execute(t *testing.T) {
	ctx := context.Background()

	rate := 4.5
	content := "とても良かった"

	t.Run("success_with_rate_and_content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(1)).
			Return(&model.DateSpotReview{ID: 1, DateSpotID: 3}, nil)
		reviewRepo.EXPECT().
			UpdateByID(ctx, uint(1), gomock.Any()).
			Return(nil)
		reviewRepo.EXPECT().
			FindByDateSpotID(ctx, uint(3)).
			Return([]*model.DateSpotReview{}, nil)

		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			FindByID(ctx, uint(3)).
			Return(&model.DateSpot{ID: 3, Name: "テストスポット"}, nil)

		interactor := usecase.NewUpdateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.UpdateDateSpotReviewInput{
			ReviewID: 1,
			Rate:     &rate,
			Content:  &content,
		})

		require.NoError(t, err)
		assert.Equal(t, uint(1), output.ReviewID)
		assert.NotNil(t, output.DateSpot)
		assert.NotNil(t, output.DateSpotReviews)
	})

	t.Run("success_with_rate_only", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(1)).
			Return(&model.DateSpotReview{ID: 1, DateSpotID: 3}, nil)
		reviewRepo.EXPECT().
			UpdateByID(ctx, uint(1), gomock.Any()).
			Return(nil)
		reviewRepo.EXPECT().
			FindByDateSpotID(ctx, uint(3)).
			Return([]*model.DateSpotReview{}, nil)

		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			FindByID(ctx, uint(3)).
			Return(&model.DateSpot{ID: 3, Name: "テストスポット"}, nil)

		interactor := usecase.NewUpdateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.UpdateDateSpotReviewInput{
			ReviewID: 1,
			Rate:     &rate,
		})

		require.NoError(t, err)
		assert.Equal(t, uint(1), output.ReviewID)
	})

	t.Run("error_validation_no_fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)

		interactor := usecase.NewUpdateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.UpdateDateSpotReviewInput{
			ReviewID: 1,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_review_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(1)).
			Return(nil, errors.New("not found"))

		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)

		interactor := usecase.NewUpdateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.UpdateDateSpotReviewInput{
			ReviewID: 1,
			Rate:     &rate,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_repository_update_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(1)).
			Return(&model.DateSpotReview{ID: 1, DateSpotID: 3}, nil)
		reviewRepo.EXPECT().
			UpdateByID(ctx, uint(1), gomock.Any()).
			Return(errors.New("db error"))

		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)

		interactor := usecase.NewUpdateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.UpdateDateSpotReviewInput{
			ReviewID: 1,
			Rate:     &rate,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
