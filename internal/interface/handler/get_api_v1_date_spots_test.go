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
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func dummyDateSpot(id uint, name string) *model.DateSpot {
	return &model.DateSpot{
		ID:       id,
		Name:     name,
		CityName: "渋谷区",
	}
}

func TestGetApiV1DateSpotsHandler(t *testing.T) {
	t.Run("success_returns_200_with_date_spots", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dateSpot := dummyDateSpot(1, "東京タワー")
		mockPort := usecasemock.NewMockGetDateSpotsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.GetDateSpotsOutput{DateSpots: []*model.DateSpot{dateSpot}}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/date_spots", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.GetApiV1DateSpots(ctx, openapi.GetApiV1DateSpotsParams{})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Len(t, resp, 1)
	})

	t.Run("success_returns_200_empty_list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetDateSpotsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.GetDateSpotsOutput{DateSpots: []*model.DateSpot{}}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/date_spots", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.GetApiV1DateSpots(ctx, openapi.GetApiV1DateSpotsParams{})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Len(t, resp, 0)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetDateSpotsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/date_spots", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.GetApiV1DateSpots(ctx, openapi.GetApiV1DateSpotsParams{})

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
