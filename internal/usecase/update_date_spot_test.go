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

func TestUpdateDateSpotInteractor_Execute(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Update(ctx, uint(10), gomock.Any()).
			Return(nil)

		interactor := usecase.NewUpdateDateSpotUsecase(dateSpotRepo)
		err := interactor.Execute(ctx, usecase.UpdateDateSpotInput{
			DateSpotID:   10,
			Name:         "更新スポット",
			GenreID:      2,
			PrefectureID: 14,
			CityName:     "新宿区",
			OpeningTime:  &now,
			ClosingTime:  &now,
		})

		require.NoError(t, err)
	})

	t.Run("error_repository_update_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		dateSpotRepo := repositorymock.NewMockDateSpotRepository(ctrl)
		dateSpotRepo.EXPECT().
			Update(ctx, uint(10), gomock.Any()).
			Return(errors.New("db error"))

		interactor := usecase.NewUpdateDateSpotUsecase(dateSpotRepo)
		err := interactor.Execute(ctx, usecase.UpdateDateSpotInput{
			DateSpotID:   10,
			Name:         "更新スポット",
			GenreID:      2,
			PrefectureID: 14,
			CityName:     "新宿区",
			OpeningTime:  &now,
			ClosingTime:  &now,
		})

		assert.Error(t, err)
	})
}
