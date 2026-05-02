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

func TestDeleteApiV1DateSpotReviewsIdHandler(t *testing.T) {
	t.Run("success_returns_200_with_reviews", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteDateSpotReviewInput{ReviewID: 10}).
			Return(&usecase.DeleteDateSpotReviewOutput{
				DateSpotReviews: []*model.DateSpotReview{},
			}, nil)

		ctx, rec := setupFormRequest(http.MethodDelete, "/api/v1/date_spot_reviews/10", url.Values{})

		h := handler.DeleteApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1DateSpotReviewsId(ctx, 10)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp, "date_spot_reviews")
		assert.Contains(t, resp, "review_average_rate")
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		ctx, _ := setupFormRequest(http.MethodDelete, "/api/v1/date_spot_reviews/10", url.Values{})

		h := handler.DeleteApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1DateSpotReviewsId(ctx, 10)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
