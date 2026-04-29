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

func TestGetApiV1GenresIdHandler(t *testing.T) {
	t.Run("success_returns_200_with_date_spots", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		genreID := 2
		mockPort := usecasemock.NewMockGetDateSpotsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetDateSpotsInput{GenreID: &genreID}).
			Return(&usecase.GetDateSpotsOutput{DateSpots: []*model.DateSpot{}}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/genres/2", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1GenresIdHandler{InputPort: mockPort}
		err := h.GetApiV1GenresId(ctx, 2)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp, "date_spots")
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		genreID := 2
		mockPort := usecasemock.NewMockGetDateSpotsInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetDateSpotsInput{GenreID: &genreID}).
			Return(nil, apperror.InternalServerError(nil))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/genres/2", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1GenresIdHandler{InputPort: mockPort}
		err := h.GetApiV1GenresId(ctx, 2)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
