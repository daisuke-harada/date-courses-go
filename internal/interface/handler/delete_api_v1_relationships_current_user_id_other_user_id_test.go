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

func TestDeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler(t *testing.T) {
	t.Run("success_returns_200_with_users", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentUwr := &model.UserWithRelations{
			User:         &model.User{ID: 1, Name: "current", Email: "current@example.com", Gender: "男性"},
			FollowerIDs:  []int{},
			FollowingIDs: []int{},
		}
		unfollowedUwr := &model.UserWithRelations{
			User:         &model.User{ID: 2, Name: "unfollowed", Email: "unfollowed@example.com", Gender: "女性"},
			FollowerIDs:  []int{},
			FollowingIDs: []int{},
		}

		mockPort := usecasemock.NewMockDeleteRelationshipInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.DeleteRelationshipInput{UserID: 1, FollowID: 2}).
			Return(&usecase.DeleteRelationshipOutput{
				Users:          []*model.UserWithRelations{currentUwr, unfollowedUwr},
				CurrentUser:    currentUwr,
				UnfollowedUser: unfollowedUwr,
			}, nil)

		ctx, rec := setupFormRequest(http.MethodDelete, "/api/v1/relationships/1/2", url.Values{})

		h := handler.DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler{InputPort: mockPort}
		err := h.DeleteApiV1RelationshipsCurrentUserIdOtherUserId(ctx, 1, 2)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp, "users")
		assert.Contains(t, resp, "current_user")
		assert.Contains(t, resp, "unfollowed_user")
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
