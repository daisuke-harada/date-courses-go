package handler_test

import (
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

func validSignupForm() url.Values {
	form := url.Values{}
	form.Set("name", "新規ユーザー")
	form.Set("email", "newuser@example.com")
	form.Set("gender", "女性")
	form.Set("password", "password123")
	form.Set("password_confirmation", "password123")
	return form
}

func TestPostApiV1SignupHandler(t *testing.T) {
	t.Run("success_returns_201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 5, Name: "新規ユーザー", Email: "newuser@example.com", Gender: model.GenderFemale}
		mockPort := usecasemock.NewMockSignupInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.SignupOutput{User: user}, nil)

		ctx, rec := setupFormRequest(http.MethodPost, "/api/v1/signup", validSignupForm())

		h := handler.PostApiV1SignupHandler{InputPort: mockPort}
		err := h.PostApiV1Signup(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("error_usecase_returns_unprocessable_entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockSignupInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.UnprocessableEntity("名前を入力してください"))

		form := validSignupForm()
		form.Del("name")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/signup", form)

		h := handler.PostApiV1SignupHandler{InputPort: mockPort}
		err := h.PostApiV1Signup(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockSignupInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/signup", validSignupForm())

		h := handler.PostApiV1SignupHandler{InputPort: mockPort}
		err := h.PostApiV1Signup(ctx)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
