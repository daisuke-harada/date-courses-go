package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repositorymock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetDateSpotsInteractor_Execute(t *testing.T) {
	t.Run("success_returns_date_spots", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		dateSpots := []*model.DateSpot{
			{ID: 1, Name: "東京タワー", CityName: "港区"},
			{ID: 2, Name: "スカイツリー", CityName: "墨田区"},
		}

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Search(ctx, gomock.Any()).
			Return(dateSpots, nil)

		interactor := usecase.NewGetDateSpotsUsecase(dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.GetDateSpotsInput{})

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Len(t, output.DateSpots, 2)
	})

	t.Run("success_with_name_filter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		name := "東京"
		dateSpots := []*model.DateSpot{{ID: 1, Name: "東京タワー", CityName: "港区"}}

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Search(ctx, gomock.Any()).
			Return(dateSpots, nil)

		interactor := usecase.NewGetDateSpotsUsecase(dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.GetDateSpotsInput{DateSpotName: &name})

		require.NoError(t, err)
		assert.Len(t, output.DateSpots, 1)
	})

	t.Run("error_repository_search_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Search(ctx, gomock.Any()).
			Return(nil, errors.New("db error"))

		interactor := usecase.NewGetDateSpotsUsecase(dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.GetDateSpotsInput{})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
