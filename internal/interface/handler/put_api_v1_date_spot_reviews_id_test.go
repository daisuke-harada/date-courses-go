package handler_test

import (
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

func validUpdateDateSpotReviewForm() url.Values {
	form := url.Values{}
	form.Set("rate", "4.5")
	form.Set("content", "更新後のレビュー")
	return form
}

func TestPutApiV1DateSpotReviewsIdHandler(t *testing.T) {
	t.Run("success_returns_200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.UpdateDateSpotReviewOutput{}, nil)

		ctx, rec := setupFormRequest(http.MethodPut, "/api/v1/date_spot_reviews/10", validUpdateDateSpotReviewForm())

		h := handler.PutApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotReviewsId(ctx, 10)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.NotFound())

		ctx, _ := setupFormRequest(http.MethodPut, "/api/v1/date_spot_reviews/99", validUpdateDateSpotReviewForm())

		h := handler.PutApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotReviewsId(ctx, 99)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateDateSpotReviewInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		ctx, _ := setupFormRequest(http.MethodPut, "/api/v1/date_spot_reviews/10", validUpdateDateSpotReviewForm())

		h := handler.PutApiV1DateSpotReviewsIdHandler{InputPort: mockPort}
		err := h.PutApiV1DateSpotReviewsId(ctx, 10)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
