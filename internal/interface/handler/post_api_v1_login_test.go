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

func TestPostApiV1LoginHandler(t *testing.T) {
	t.Run("success_returns_200_with_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "testuser", Email: "test@example.com", Gender: model.GenderMale}
		mockPort := usecasemock.NewMockLoginInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.LoginOutput{User: user, Token: "test-jwt-token"}, nil)

		form := url.Values{}
		form.Set("name", "testuser")
		form.Set("password", "password123")
		ctx, rec := setupFormRequest(http.MethodPost, "/api/v1/login", form)

		h := handler.PostApiV1LoginHandler{InputPort: mockPort}
		err := h.PostApiV1Login(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "test-jwt-token", resp["token"])
	})

	t.Run("error_unauthorized_invalid_credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockLoginInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.Unauthorized("名前またはパスワードが正しくありません"))

		form := url.Values{}
		form.Set("name", "wronguser")
		form.Set("password", "wrongpass")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/login", form)

		h := handler.PostApiV1LoginHandler{InputPort: mockPort}
		err := h.PostApiV1Login(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockLoginInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		form := url.Values{}
		form.Set("name", "testuser")
		form.Set("password", "password123")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/login", form)

		h := handler.PostApiV1LoginHandler{InputPort: mockPort}
		err := h.PostApiV1Login(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
