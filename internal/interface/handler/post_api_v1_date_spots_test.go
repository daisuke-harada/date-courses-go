package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func validDateSpotForm() url.Values {
	form := url.Values{}
	form.Set("name", "東京タワー")
	form.Set("genre_id", "1")
	form.Set("prefecture_id", "13")
	form.Set("city_name", "港区")
	form.Set("opening_time", "2024-01-01T09:00:00Z")
	form.Set("closing_time", "2024-01-01T21:00:00Z")
	return form
}

func TestPostApiV1DateSpotsHandler(t *testing.T) {
	t.Run("success_returns_201_with_date_spot_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.CreateDateSpotOutput{DateSpotID: 42}, nil)

		ctx, rec := setupFormRequest(http.MethodPost, "/api/v1/date_spots", validDateSpotForm())

		h := handler.PostApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpots(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, float64(42), resp["date_spot_id"])
	})

	t.Run("error_missing_name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotInputPort(ctrl)

		form := validDateSpotForm()
		form.Del("name")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spots", form)

		h := handler.PostApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpots(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_missing_genre_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotInputPort(ctrl)

		form := validDateSpotForm()
		form.Del("genre_id")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spots", form)

		h := handler.PostApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpots(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_missing_city_name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotInputPort(ctrl)

		form := validDateSpotForm()
		form.Del("city_name")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spots", form)

		h := handler.PostApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpots(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_invalid_opening_time_format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotInputPort(ctrl)

		form := validDateSpotForm()
		form.Set("opening_time", "not-a-valid-time")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spots", form)

		h := handler.PostApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpots(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spots", validDateSpotForm())

		h := handler.PostApiV1DateSpotsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpots(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
