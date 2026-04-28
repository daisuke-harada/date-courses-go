package handler_test

import (
	"encoding/json"
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

func TestGetApiV1TopHandler(t *testing.T) {
	t.Run("success_returns_200_with_required_fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetDateSpotsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetDateSpotsInput{}).
			Return(&usecase.GetDateSpotsOutput{DateSpots: []*model.DateSpot{}}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/top", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1TopHandler{InputPort: mockPort}
		err := h.GetApiV1Top(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp, "date_spots")
		assert.Contains(t, resp, "areas")
		assert.Contains(t, resp, "genres")
		assert.Contains(t, resp, "main_genres")
		assert.Contains(t, resp, "main_prefecture")
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetDateSpotsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetDateSpotsInput{}).
			Return(nil, apperror.InternalServerError(nil))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/top", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1TopHandler{InputPort: mockPort}
		err := h.GetApiV1Top(ctx)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
