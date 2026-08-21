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

func TestGetCourseInteractor_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		courseRepo := repomock.NewMockCourseRepository(ctrl)
		courseRepo.EXPECT().
			FindByID(ctx, uint(1), gomock.Any()).
			Return(&model.Course{ID: 1, UserID: 2, TravelMode: "DRIVING", Authority: "公開"}, nil)

		interactor := usecase.NewGetCourseUsecase(courseRepo)
		output, err := interactor.Execute(ctx, usecase.GetCourseInput{CourseID: 1})

		require.NoError(t, err)
		assert.Equal(t, uint(1), output.Course.ID)
		assert.Equal(t, uint(2), output.Course.UserID)
	})

	// 閲覧者の ID がそのまま repository に渡り、SQL 側で可視性が絞られる
	t.Run("passes_viewer_id_to_repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		courseRepo := repomock.NewMockCourseRepository(ctrl)
		courseRepo.EXPECT().
			FindByID(ctx, uint(1), uint(7)).
			Return(&model.Course{ID: 1, UserID: 7, Authority: model.CourseAuthorityPrivate}, nil)

		interactor := usecase.NewGetCourseUsecase(courseRepo)
		output, err := interactor.Execute(ctx, usecase.GetCourseInput{CourseID: 1, ViewerID: 7})

		require.NoError(t, err)
		assert.Equal(t, uint(1), output.Course.ID)
	})

	// 見えないコースは repository が返さないので 404 になる
	t.Run("error_not_found_when_invisible", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		courseRepo := repomock.NewMockCourseRepository(ctrl)
		courseRepo.EXPECT().
			FindByID(ctx, uint(1), uint(0)).
			Return(nil, errors.New("record not found"))

		interactor := usecase.NewGetCourseUsecase(courseRepo)
		output, err := interactor.Execute(ctx, usecase.GetCourseInput{CourseID: 1, ViewerID: 0})

		assert.Nil(t, output)
		require.Error(t, err)
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		courseRepo := repomock.NewMockCourseRepository(ctrl)
		courseRepo.EXPECT().
			FindByID(ctx, uint(999), gomock.Any()).
			Return(nil, errors.New("record not found"))

		interactor := usecase.NewGetCourseUsecase(courseRepo)
		output, err := interactor.Execute(ctx, usecase.GetCourseInput{CourseID: 999})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
