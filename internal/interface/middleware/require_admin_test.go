package middleware_test

import (
	"net/http"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAdmin(t *testing.T) {
	t.Run("returns_user_when_admin", func(t *testing.T) {
		ctx := newContext()
		middleware.SetCurrentUser(ctx, &model.User{ID: 1, Name: "admin", Admin: true})

		user, err := middleware.RequireAdmin(ctx)

		require.NoError(t, err)
		assert.Equal(t, uint(1), user.ID)
	})

	// 一般ユーザーがデートスポットを操作できないことを保証する
	t.Run("error_forbidden_when_not_admin", func(t *testing.T) {
		ctx := newContext()
		middleware.SetCurrentUser(ctx, &model.User{ID: 2, Name: "alice", Admin: false})

		_, err := middleware.RequireAdmin(ctx)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, statusCode)
	})

	t.Run("error_unauthorized_when_anonymous", func(t *testing.T) {
		_, err := middleware.RequireAdmin(newContext())

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
	})
}
