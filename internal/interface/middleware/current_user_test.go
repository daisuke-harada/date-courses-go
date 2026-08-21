package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/courses", nil)
	return e.NewContext(req, httptest.NewRecorder())
}

func TestCurrentUser(t *testing.T) {
	t.Run("returns_user_when_authenticated", func(t *testing.T) {
		ctx := newContext()
		middleware.SetCurrentUser(ctx, &model.User{ID: 7, Name: "alice"})

		user := middleware.CurrentUser(ctx)

		require.NotNil(t, user)
		assert.Equal(t, uint(7), user.ID)
	})

	t.Run("returns_nil_when_anonymous", func(t *testing.T) {
		assert.Nil(t, middleware.CurrentUser(newContext()))
	})
}

func TestRequireCurrentUser(t *testing.T) {
	t.Run("returns_user_when_authenticated", func(t *testing.T) {
		ctx := newContext()
		middleware.SetCurrentUser(ctx, &model.User{ID: 7, Name: "alice"})

		user, err := middleware.RequireCurrentUser(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint(7), user.ID)
	})

	t.Run("error_unauthorized_when_anonymous", func(t *testing.T) {
		_, err := middleware.RequireCurrentUser(newContext())

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
	})
}
