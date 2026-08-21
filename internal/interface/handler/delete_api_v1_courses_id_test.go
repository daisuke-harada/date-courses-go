package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupDeleteCourseRequest() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/courses/1", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestDeleteApiV1CoursesIdHandler(t *testing.T) {
	t.Run("success_returns_204", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteCourseInput{CourseID: 1, OperatorID: 10}).
			Return(nil)

		ctx, rec := setupDeleteCourseRequest()
		middleware.SetCurrentUser(ctx, &model.User{ID: 10, Name: "alice"})

		h := handler.DeleteApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1CoursesId(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("error_unauthorized_without_current_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteCourseInputPort(ctrl)

		ctx, _ := setupDeleteCourseRequest()

		h := handler.DeleteApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1CoursesId(ctx, 1)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
	})

	t.Run("error_usecase_returns_forbidden", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteCourseInput{CourseID: 1, OperatorID: 99}).
			Return(apperror.Forbidden("他のユーザーのデートコースは削除できません"))

		ctx, _ := setupDeleteCourseRequest()
		middleware.SetCurrentUser(ctx, &model.User{ID: 99, Name: "bob"})

		h := handler.DeleteApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1CoursesId(ctx, 1)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteCourseInput{CourseID: 1, OperatorID: 10}).
			Return(apperror.InternalServerError(nil))

		ctx, _ := setupDeleteCourseRequest()
		middleware.SetCurrentUser(ctx, &model.User{ID: 10, Name: "alice"})

		h := handler.DeleteApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1CoursesId(ctx, 1)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
