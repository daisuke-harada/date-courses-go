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

func TestGetApiV1DateSpotsIdHandler(t *testing.T) {
	t.Run("success_returns_200_with_date_spot_reviews_and_average_rate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dateSpot := dummyDateSpot(1, "東京タワー")
		rate := 4.0
		content := "良かった"
		userName := "田中"
		gender := model.GenderMale
		reviews := []*model.DateSpotReview{
			{
				ID:         1,
				DateSpotID: 1,
				UserID:     1,
				Rate:       &rate,
				Content:    &content,
				User:       &model.User{ID: 1, Name: userName, Gender: gender},
			},
		}

		mockPort := usecasemock.NewMockGetDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetDateSpotInput{ID: 1}).
			Return(&usecase.GetDateSpotOutput{DateSpot: dateSpot, DateSpotReviews: reviews}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/date_spots/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.GetApiV1DateSpotsId(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp, "date_spot")
		assert.Contains(t, resp, "review_average_rate")
		assert.Contains(t, resp, "date_spot_reviews")

		reviewsList, ok := resp["date_spot_reviews"].([]interface{})
		require.True(t, ok)
		assert.Len(t, reviewsList, 1)
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.NotFound())

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/date_spots/999", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.GetApiV1DateSpotsId(ctx, 999)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/date_spots/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1DateSpotsIdHandler{InputPort: mockPort}
		err := h.GetApiV1DateSpotsId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
