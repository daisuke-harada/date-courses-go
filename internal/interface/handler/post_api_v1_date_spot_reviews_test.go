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

func validDateSpotReviewForm() url.Values {
	form := url.Values{}
	form.Set("user_id", "1")
	form.Set("date_spot_id", "2")
	form.Set("rate", "4.5")
	form.Set("content", "とても良かった")
	return form
}

func TestPostApiV1DateSpotReviewsHandler(t *testing.T) {
	t.Run("success_returns_201_with_reviews", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.CreateDateSpotReviewOutput{
				ReviewID:        10,
				DateSpotReviews: []*model.DateSpotReview{},
			}, nil)

		ctx, rec := setupFormRequest(http.MethodPost, "/api/v1/date_spot_reviews", validDateSpotReviewForm())

		h := handler.PostApiV1DateSpotReviewsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpotReviews(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp, "date_spot_reviews")
		assert.Contains(t, resp, "review_average_rate")
	})

	t.Run("error_bad_request_invalid_user_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotReviewInputPort(ctrl)

		form := validDateSpotReviewForm()
		form.Set("user_id", "abc")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spot_reviews", form)

		h := handler.PostApiV1DateSpotReviewsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpotReviews(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("error_bad_request_invalid_date_spot_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotReviewInputPort(ctrl)

		form := validDateSpotReviewForm()
		form.Set("date_spot_id", "xyz")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spot_reviews", form)

		h := handler.PostApiV1DateSpotReviewsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpotReviews(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("error_usecase_returns_unprocessable_entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.UnprocessableEntity("ユーザーIDを入力してください"))

		form := validDateSpotReviewForm()
		form.Set("user_id", "0")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spot_reviews", form)

		h := handler.PostApiV1DateSpotReviewsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpotReviews(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spot_reviews", validDateSpotReviewForm())

		h := handler.PostApiV1DateSpotReviewsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpotReviews(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
