package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repositorymock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateDateSpotInteractor_Execute(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, ds *model.DateSpot) error {
				ds.ID = 10
				return nil
			})

		interactor := usecase.NewCreateDateSpotUsecase(dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.CreateDateSpotInput{
			Name:         "テストスポット",
			GenreID:      1,
			PrefectureID: 13,
			CityName:     "渋谷区",
			OpeningTime:  &now,
			ClosingTime:  &now,
		})

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, uint(10), output.DateSpotID)
	})

	t.Run("error_repository_create_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("db error"))

		interactor := usecase.NewCreateDateSpotUsecase(dateSpotRepo)
		output, err := interactor.Execute(ctx, usecase.CreateDateSpotInput{
			Name:         "テストスポット",
			GenreID:      1,
			PrefectureID: 13,
			CityName:     "渋谷区",
			OpeningTime:  &now,
			ClosingTime:  &now,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
