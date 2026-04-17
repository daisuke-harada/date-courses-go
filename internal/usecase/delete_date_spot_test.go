package usecase_test

import (
	"context"
	"errors"
	"testing"

	repositorymock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteDateSpotInteractor_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Delete(ctx, uint(10)).
			Return(nil)

		interactor := usecase.NewDeleteDateSpotUsecase(dateSpotRepo)
		err := interactor.Execute(ctx, usecase.DeleteDateSpotInput{DateSpotID: 10})

		require.NoError(t, err)
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Delete(ctx, uint(10)).
			Return(errors.New("not found"))

		interactor := usecase.NewDeleteDateSpotUsecase(dateSpotRepo)
		err := interactor.Execute(ctx, usecase.DeleteDateSpotInput{DateSpotID: 10})

		assert.Error(t, err)
	})
}
