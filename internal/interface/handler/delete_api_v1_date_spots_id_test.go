package handler_test

import (
	"errors"
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

func TestDeleteApiV1DateSpotsIdHandler(t *testing.T) {
	// デートスポットの削除は管理者だけ。一般ユーザーは拒否される
	t.Run("error_forbidden_for_non_admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteDateSpotInputPort(ctrl)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/date_spots/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		middleware.SetCurrentUser(ctx, &model.User{ID: 2, Name: "alice", Admin: false})

		h := handler.DeleteApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1DateSpotsId(ctx, 1)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, statusCode)
	})

	t.Run("error_unauthorized_without_current_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteDateSpotInputPort(ctrl)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/date_spots/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.DeleteApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1DateSpotsId(ctx, 1)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
	})

	t.Run("success_returns_204", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteDateSpotInput{DateSpotID: 10}).
			Return(nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/date_spots/10", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("10")

		middleware.SetCurrentUser(ctx, &model.User{ID: 1, Name: "admin", Admin: true})
		h := handler.DeleteApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1DateSpotsId(ctx, 10)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteDateSpotInput{DateSpotID: 10}).
			Return(apperror.InternalServerError(errors.New("db error")))

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/date_spots/10", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("10")

		middleware.SetCurrentUser(ctx, &model.User{ID: 1, Name: "admin", Admin: true})
		h := handler.DeleteApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1DateSpotsId(ctx, 10)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
