package handler_test

import (
	"errors"
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

func TestDeleteApiV1UsersIdHandler(t *testing.T) {
	t.Run("success_returns_204", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteUserInput{ID: 1}).
			Return(nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.DeleteApiV1UsersIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1UsersId(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(apperror.NotFound())

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/999", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.DeleteApiV1UsersIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1UsersId(ctx, 999)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(apperror.InternalServerError(errors.New("db error")))

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.DeleteApiV1UsersIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1UsersId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
