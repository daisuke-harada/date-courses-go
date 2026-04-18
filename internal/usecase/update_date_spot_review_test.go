package usecase_test

import (
	"context"
	"errors"
	"testing"

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
	content := "更新後のレビュー"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(10)).
			Return(&model.DateSpotReview{ID: 10, UserID: 1, DateSpotID: 2}, nil)
		reviewRepo.EXPECT().
			Update(ctx, gomock.Any()).
			Return(nil)

		interactor := usecase.NewUpdateDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.UpdateDateSpotReviewInput{
			ReviewID: 10,
			Rate:     &rate,
			Content:  &content,
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(99)).
			Return(nil, errors.New("record not found"))

		interactor := usecase.NewUpdateDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.UpdateDateSpotReviewInput{ReviewID: 99})

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("error_update_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			FindByID(ctx, uint(10)).
			Return(&model.DateSpotReview{ID: 10}, nil)
		reviewRepo.EXPECT().
			Update(ctx, gomock.Any()).
			Return(errors.New("db error"))

		interactor := usecase.NewUpdateDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.UpdateDateSpotReviewInput{ReviewID: 10, Rate: &rate})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
