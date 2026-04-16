package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

// ─── テストヘルパー ───────────────────────────────────────────────────────────

func setupFormRequest(method, path string, form url.Values) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	return ctx, rec
}

func dummyUserWithRelations(id uint, name string) *model.UserWithRelations {
	return &model.UserWithRelations{
		User:         &model.User{ID: id, Name: name, Email: name + "@example.com", Gender: model.GenderMale},
		FollowerIDs:  []int{},
		FollowingIDs: []int{},
		Courses:      []*model.Course{},
		Reviews:      []*model.DateSpotReview{},
	}
}

// ─── テスト ───────────────────────────────────────────────────────────────────

func TestPostApiV1RelationshipsHandler(t *testing.T) {
	t.Run("success_returns_201_with_correct_json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		currentUser := dummyUserWithRelations(1, "alice")
		followedUser := dummyUserWithRelations(2, "bob")

		mockPort := usecasemock.NewMockCreateRelationshipInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.CreateRelationshipInput{
				CurrentUserID:  1,
				FollowedUserID: 2,
			}).
			Return(&usecase.CreateRelationshipOutput{
				Users:        []*model.UserWithRelations{currentUser, followedUser},
				CurrentUser:  currentUser,
				FollowedUser: followedUser,
			}, nil)

		form := url.Values{}
		form.Set("current_user_id", "1")
		form.Set("followed_user_id", "2")
		ctx, rec := setupFormRequest(http.MethodPost, "/api/v1/relationships", form)

		h := handler.PostApiV1RelationshipsHandler{InputPort: mockPort}
		err := h.PostApiV1Relationships(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp, "users")
		assert.Contains(t, resp, "current_user")
		assert.Contains(t, resp, "followed_user")
		users := resp["users"].([]interface{})
		assert.Equal(t, 2, len(users))
	})

	t.Run("error_invalid_current_user_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateRelationshipInputPort(ctrl)
		// Execute は呼ばれないので EXPECT 不要

		form := url.Values{}
		form.Set("current_user_id", "abc")
		form.Set("followed_user_id", "2")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/relationships", form)

		h := handler.PostApiV1RelationshipsHandler{InputPort: mockPort}
		err := h.PostApiV1Relationships(ctx)

		assert.Error(t, err)
		_, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok, "apperror 型のエラーであること")
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateRelationshipInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.CreateRelationshipInput{
				CurrentUserID:  1,
				FollowedUserID: 1,
			}).
			Return(nil, apperror.UnprocessableEntity("自分自身をフォローすることはできません"))

		form := url.Values{}
		form.Set("current_user_id", "1")
		form.Set("followed_user_id", "1")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/relationships", form)

		h := handler.PostApiV1RelationshipsHandler{InputPort: mockPort}
		err := h.PostApiV1Relationships(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})
}
