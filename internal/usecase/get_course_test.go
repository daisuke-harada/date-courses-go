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
			FindByID(ctx, uint(1)).
			Return(&model.Course{ID: 1, UserID: 2, TravelMode: "DRIVING", Authority: "公開"}, nil)

		interactor := usecase.NewGetCourseUsecase(courseRepo)
		output, err := interactor.Execute(ctx, usecase.GetCourseInput{CourseID: 1})

		require.NoError(t, err)
		assert.Equal(t, uint(1), output.Course.ID)
		assert.Equal(t, uint(2), output.Course.UserID)
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		courseRepo := repomock.NewMockCourseRepository(ctrl)
		courseRepo.EXPECT().
			FindByID(ctx, uint(999)).
			Return(nil, errors.New("record not found"))

		interactor := usecase.NewGetCourseUsecase(courseRepo)
		output, err := interactor.Execute(ctx, usecase.GetCourseInput{CourseID: 999})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
