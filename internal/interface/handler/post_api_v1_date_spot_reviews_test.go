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

func validDateSpotReviewForm() url.Values {
	form := url.Values{}
	form.Set("user_id", "1")
	form.Set("date_spot_id", "2")
	form.Set("rate", "4.5")
	form.Set("content", "とても良かった")
	return form
}

func TestPostApiV1DateSpotReviewsHandler(t *testing.T) {
	t.Run("success_returns_201_with_review_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.CreateDateSpotReviewOutput{ReviewID: 10}, nil)

		ctx, rec := setupFormRequest(http.MethodPost, "/api/v1/date_spot_reviews", validDateSpotReviewForm())

		h := handler.PostApiV1DateSpotReviewsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpotReviews(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, float64(10), resp["review_id"])
	})

	t.Run("error_missing_user_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotReviewInputPort(ctrl)

		form := validDateSpotReviewForm()
		form.Del("user_id")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/date_spot_reviews", form)

		h := handler.PostApiV1DateSpotReviewsHandler{InputPort: mockPort}
		err := h.PostApiV1DateSpotReviews(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_missing_date_spot_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateDateSpotReviewInputPort(ctrl)

		form := validDateSpotReviewForm()
		form.Del("date_spot_id")
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
