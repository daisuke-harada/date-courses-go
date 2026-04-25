package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
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
			DeleteByID(gomock.Any(), uint(1)).
			Return(nil)

		uc := usecase.NewDeleteCourseUsecase(mockCourseRepo)
		require.NoError(t, uc.Execute(context.Background(), usecase.DeleteCourseInput{CourseID: 1}))
	})

	t.Run("error_repository_delete_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockCourseRepo.EXPECT().
			DeleteByID(gomock.Any(), uint(1)).
			Return(errors.New("db error"))

		uc := usecase.NewDeleteCourseUsecase(mockCourseRepo)
		err := uc.Execute(context.Background(), usecase.DeleteCourseInput{CourseID: 1})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
