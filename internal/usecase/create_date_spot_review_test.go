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

func TestCreateDateSpotReviewInteractor_Execute(t *testing.T) {
	ctx := context.Background()

	rate := 4.5
	content := "とても良かった"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, r *model.DateSpotReview) error {
				r.ID = 10
				return nil
			})
		reviewRepo.EXPECT().
			FindByDateSpotID(ctx, uint(2)).
			Return([]*model.DateSpotReview{}, nil)

		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			FindByID(ctx, uint(2)).
			Return(&model.DateSpot{ID: 2, Name: "テストスポット"}, nil)

		interactor := usecase.NewCreateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.CreateDateSpotReviewInput{
			UserID:     1,
			DateSpotID: 2,
			Rate:       &rate,
			Content:    &content,
		})

		require.NoError(t, err)
		assert.Equal(t, uint(10), output.ReviewID)
		assert.NotNil(t, output.DateSpot)
		assert.NotNil(t, output.DateSpotReviews)
	})

	t.Run("error_validation_missing_user_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)

		interactor := usecase.NewCreateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.CreateDateSpotReviewInput{
			UserID:     0,
			DateSpotID: 2,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_validation_missing_date_spot_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)

		interactor := usecase.NewCreateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.CreateDateSpotReviewInput{
			UserID:     1,
			DateSpotID: 0,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_repository_create_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("db error"))

		dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)

		interactor := usecase.NewCreateDateSpotReviewUsecase(reviewRepo, dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.CreateDateSpotReviewInput{
			UserID:     1,
			DateSpotID: 2,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
