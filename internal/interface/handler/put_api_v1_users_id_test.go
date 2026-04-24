package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func validUpdateUserForm() url.Values {
	form := url.Values{}
	form.Set("name", "テストユーザー")
	form.Set("email", "test@example.com")
	form.Set("gender", "男性")
	return form
}

func TestPutApiV1UsersIdHandler(t *testing.T) {
	t.Run("success_returns_200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := dummyUserWithRelations(1, "テストユーザー")
		mockPort := usecasemock.NewMockUpdateUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&usecase.UpdateUserOutput{UserWithRelations: user}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", strings.NewReader(validUpdateUserForm().Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		h := handler.PutApiV1UsersIdHandler{InputPort: mockPort}
		err := h.PutApiV1UsersId(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_usecase_returns_unprocessable_entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.UnprocessableEntity("名前を入力してください"))

		form := validUpdateUserForm()
		form.Del("name")
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		h := handler.PutApiV1UsersIdHandler{InputPort: mockPort}
		err := h.PutApiV1UsersId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_usecase_returns_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.NotFound())

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/999", strings.NewReader(validUpdateUserForm().Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("999")

		h := handler.PutApiV1UsersIdHandler{InputPort: mockPort}
		err := h.PutApiV1UsersId(ctx, 999)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockUpdateUserInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", strings.NewReader(validUpdateUserForm().Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		h := handler.PutApiV1UsersIdHandler{InputPort: mockPort}
		err := h.PutApiV1UsersId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
