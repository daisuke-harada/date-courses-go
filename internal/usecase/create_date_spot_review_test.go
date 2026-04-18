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

		interactor := usecase.NewCreateDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.CreateDateSpotReviewInput{
			UserID:     1,
			DateSpotID: 2,
			Rate:       &rate,
			Content:    &content,
		})

		require.NoError(t, err)
		assert.Equal(t, uint(10), output.ReviewID)
	})

	t.Run("error_repository_create_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("db error"))

		interactor := usecase.NewCreateDateSpotReviewUsecase(reviewRepo)
		output, err := interactor.Execute(ctx, usecase.CreateDateSpotReviewInput{
			UserID:     1,
			DateSpotID: 2,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
