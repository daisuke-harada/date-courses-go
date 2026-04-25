package usecase_test

import (
	"context"
	"errors"
	"testing"

	repomock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteDateSpotReviewInteractor_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			DeleteByID(ctx, uint(10)).
			Return(nil)

		interactor := usecase.NewDeleteDateSpotReviewUsecase(reviewRepo)
		err := interactor.Execute(ctx, usecase.DeleteDateSpotReviewInput{ReviewID: 10})

		require.NoError(t, err)
	})

	t.Run("error_repository_delete_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := repomock.NewMockDateSpotReviewRepository(ctrl)
		reviewRepo.EXPECT().
			DeleteByID(ctx, uint(10)).
			Return(errors.New("db error"))

		interactor := usecase.NewDeleteDateSpotReviewUsecase(reviewRepo)
		err := interactor.Execute(ctx, usecase.DeleteDateSpotReviewInput{ReviewID: 10})

		assert.Error(t, err)
	})
}
