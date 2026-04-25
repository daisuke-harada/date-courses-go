package handler_test

import (
	"encoding/json"
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

func TestGetApiV1UsersIdHandler(t *testing.T) {
	t.Run("success_returns_200_with_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := dummyUserWithRelations(1, "ターゲットユーザー")
		mockPort := usecasemock.NewMockGetUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetUserInput{ID: 1}).
			Return(&usecase.GetUserOutput{UserWithRelations: user}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1UsersIdHandler{InputPort: mockPort}
		err := h.GetApiV1UsersId(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, float64(1), resp["id"])
		assert.Equal(t, "ターゲットユーザー", resp["name"])
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.NotFound())

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/999", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1UsersIdHandler{InputPort: mockPort}
		err := h.GetApiV1UsersId(ctx, 999)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1UsersIdHandler{InputPort: mockPort}
		err := h.GetApiV1UsersId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
