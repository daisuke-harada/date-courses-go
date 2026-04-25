package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteApiV1CoursesIdHandler(t *testing.T) {
	t.Run("success_returns_204", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteCourseInput{CourseID: 1}).
			Return(nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/courses/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.DeleteApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1CoursesId(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteCourseInput{CourseID: 1}).
			Return(apperror.InternalServerError(nil))

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/courses/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.DeleteApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1CoursesId(ctx, 1)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
