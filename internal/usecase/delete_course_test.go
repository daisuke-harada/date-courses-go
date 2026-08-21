package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repomock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteCourseUsecase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockCourseRepo.EXPECT().
			FindByID(gomock.Any(), uint(1)).
			Return(&model.Course{ID: 1, UserID: 10}, nil)
		mockCourseRepo.EXPECT().
			DeleteByID(gomock.Any(), uint(1)).
			Return(nil)

		uc := usecase.NewDeleteCourseUsecase(mockCourseRepo)
		require.NoError(t, uc.Execute(context.Background(), usecase.DeleteCourseInput{CourseID: 1, OperatorID: 10}))
	})

	// 他人のコースを削除できないことを保証する
	t.Run("error_forbidden_when_operator_is_not_owner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockCourseRepo.EXPECT().
			FindByID(gomock.Any(), uint(1)).
			Return(&model.Course{ID: 1, UserID: 10}, nil)

		uc := usecase.NewDeleteCourseUsecase(mockCourseRepo)
		err := uc.Execute(context.Background(), usecase.DeleteCourseInput{CourseID: 1, OperatorID: 99})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, statusCode)
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockCourseRepo.EXPECT().
			FindByID(gomock.Any(), uint(1)).
			Return(nil, errors.New("record not found"))

		uc := usecase.NewDeleteCourseUsecase(mockCourseRepo)
		err := uc.Execute(context.Background(), usecase.DeleteCourseInput{CourseID: 1, OperatorID: 10})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_repository_delete_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockCourseRepo.EXPECT().
			FindByID(gomock.Any(), uint(1)).
			Return(&model.Course{ID: 1, UserID: 10}, nil)
		mockCourseRepo.EXPECT().
			DeleteByID(gomock.Any(), uint(1)).
			Return(errors.New("db error"))

		uc := usecase.NewDeleteCourseUsecase(mockCourseRepo)
		err := uc.Execute(context.Background(), usecase.DeleteCourseInput{CourseID: 1, OperatorID: 10})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
