package usecase_test

import (
	"context"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetCoursesInteractor_Execute(t *testing.T) {
	t.Run("success_returns_courses", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockCourseRepository(ctrl)
		courses := []*model.Course{
			{ID: 1, UserID: 1, TravelMode: "car", Authority: "public"},
			{ID: 2, UserID: 2, TravelMode: "walk", Authority: "public"},
		}
		mockRepo.EXPECT().
			FindAll(gomock.Any()).
			Return(courses, nil)

		interactor := usecase.NewGetCoursesUsecase(mockRepo)
		output, err := interactor.Execute(context.Background(), usecase.GetCoursesInput{})

		require.NoError(t, err)
		assert.Equal(t, 2, len(output.Courses))
		assert.Equal(t, uint(1), output.Courses[0].ID)
		assert.Equal(t, uint(2), output.Courses[1].ID)
	})

	t.Run("success_returns_empty_list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockCourseRepository(ctrl)
		mockRepo.EXPECT().
			FindAll(gomock.Any()).
			Return([]*model.Course{}, nil)

		interactor := usecase.NewGetCoursesUsecase(mockRepo)
		output, err := interactor.Execute(context.Background(), usecase.GetCoursesInput{})

		require.NoError(t, err)
		assert.Equal(t, 0, len(output.Courses))
	})

	t.Run("error_repository_failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockCourseRepository(ctrl)
		mockRepo.EXPECT().
			FindAll(gomock.Any()).
			Return(nil, apperror.InternalServerError(nil))

		interactor := usecase.NewGetCoursesUsecase(mockRepo)
		output, err := interactor.Execute(context.Background(), usecase.GetCoursesInput{})

		assert.Nil(t, output)
		assert.Error(t, err)
	})
}
