package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPutApiV1DateSpotReviewsIdHandler(t *testing.T) {
	t.Run("success_returns_200_with_date_spot", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.UpdateDateSpotReviewOutput{
				ReviewID:        1,
				DateSpot:        &model.DateSpot{ID: 3, Name: "テストスポット"},
				DateSpotReviews: []*model.DateSpotReview{},
			}, nil)

		form := url.Values{}
		form.Set("rate", "4.5")
		form.Set("content", "良かった")
		ctx, rec := setupFormRequest(http.MethodPut, "/api/v1/date_spot_reviews/1", form)

		h := handler.PutApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotReviewsId(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp, "date_spot")
		assert.Contains(t, resp, "date_spot_reviews")
		assert.Contains(t, resp, "review_average_rate")
	})

	t.Run("error_bad_request_invalid_rate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotReviewInputPort(ctrl)

		form := url.Values{}
		form.Set("rate", "not-a-number")
		ctx, _ := setupFormRequest(http.MethodPut, "/api/v1/date_spot_reviews/1", form)

		h := handler.PutApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotReviewsId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("error_usecase_returns_unprocessable_entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.UnprocessableEntity("rate または content のいずれかを入力してください"))

		form := url.Values{}
		ctx, _ := setupFormRequest(http.MethodPut, "/api/v1/date_spot_reviews/1", form)

		h := handler.PutApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotReviewsId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		form := url.Values{}
		form.Set("rate", "4.5")
		ctx, _ := setupFormRequest(http.MethodPut, "/api/v1/date_spot_reviews/1", form)

		h := handler.PutApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotReviewsId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
