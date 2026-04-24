package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetApiV1UsersUserIdFollowingsHandler(t *testing.T) {
	t.Run("success_returns_200_with_followings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		following := dummyUserWithRelations(2, "フォロー先ユーザー")
		mockPort := usecasemock.NewMockGetUserFollowingsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetUserFollowingsInput{UserID: 1}).
			Return(&usecase.GetUserFollowingsOutput{Users: []*model.UserWithRelations{following}}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/followings", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1UsersUserIdFollowingsHandler{InputPort: mockPort}
		err := h.GetApiV1UsersUserIdFollowings(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Len(t, resp, 1)
		assert.Equal(t, float64(2), resp[0]["id"])
	})

	t.Run("success_returns_200_empty_list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetUserFollowingsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.GetUserFollowingsOutput{Users: []*model.UserWithRelations{}}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/followings", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1UsersUserIdFollowingsHandler{InputPort: mockPort}
		err := h.GetApiV1UsersUserIdFollowings(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Len(t, resp, 0)
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetUserFollowingsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.NotFound())

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/999/followings", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1UsersUserIdFollowingsHandler{InputPort: mockPort}
		err := h.GetApiV1UsersUserIdFollowings(ctx, 999)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetUserFollowingsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/followings", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1UsersUserIdFollowingsHandler{InputPort: mockPort}
		err := h.GetApiV1UsersUserIdFollowings(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
