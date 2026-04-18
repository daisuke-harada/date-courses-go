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

func TestDeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler(t *testing.T) {
	t.Run("success_returns_200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteRelationshipInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteRelationshipInput{UserID: 1, FollowID: 2}).
			Return(&usecase.DeleteRelationshipOutput{}, nil)

		ctx, rec := setupFormRequest(http.MethodDelete, "/api/v1/relationships/1/2", url.Values{})

		h := handler.DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1RelationshipsCurrentUserIdOtherUserId(ctx, 1, 2)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockDeleteRelationshipInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		ctx, _ := setupFormRequest(http.MethodDelete, "/api/v1/relationships/1/2", url.Values{})

		h := handler.DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1RelationshipsCurrentUserIdOtherUserId(ctx, 1, 2)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
