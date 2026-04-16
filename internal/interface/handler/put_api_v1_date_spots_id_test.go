package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func validUpdateDateSpotForm() url.Values {
	form := url.Values{}
	form.Set("name", "更新東京タワー")
	form.Set("genre_id", "2")
	form.Set("prefecture_id", "13")
	form.Set("city_name", "港区")
	form.Set("opening_time", "2024-01-01T10:00:00Z")
	form.Set("closing_time", "2024-01-01T20:00:00Z")
	return form
}

func TestPutApiV1DateSpotsIdHandler(t *testing.T) {
	t.Run("success_returns_200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/date_spots/10", strings.NewReader(validUpdateDateSpotForm().Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("10")

		h := handler.PutApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotsId(ctx, 10)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_missing_name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotInputPort(ctrl)

		e := echo.New()
		form := validUpdateDateSpotForm()
		form.Del("name")
		req := httptest.NewRequest(http.MethodPut, "/api/v1/date_spots/10", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("10")

		h := handler.PutApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotsId(ctx, 10)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(apperror.InternalServerError(errors.New("db error")))

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/date_spots/10", strings.NewReader(validUpdateDateSpotForm().Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("10")

		h := handler.PutApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotsId(ctx, 10)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
