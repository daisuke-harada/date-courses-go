package usecase_test

import (
	"context"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repositorymock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetDateSpotInteractor_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		dateSpot := &model.DateSpot{ID: 1, Name: "東京タワー", CityName: "港区"}

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			FindByID(ctx, uint(1)).
			Return(dateSpot, nil)

		interactor := usecase.NewGetDateSpotUsecase(dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.GetDateSpotInput{ID: 1})

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, dateSpot, output.DateSpot)
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			FindByID(ctx, uint(999)).
			Return(nil, apperror.NotFound())

		interactor := usecase.NewGetDateSpotUsecase(dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.GetDateSpotInput{ID: 999})

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 404, statusCode)
	})
}
