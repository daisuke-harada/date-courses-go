package jwt_test

import (
	"testing"
	"time"

	jwtpkg "github.com/daisuke-harada/date-courses-go/internal/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

func TestEncode(t *testing.T) {
	t.Run("success_returns_non_empty_token", func(t *testing.T) {
		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestDecode(t *testing.T) {
	t.Run("success_returns_user_id", func(t *testing.T) {
		token, err := jwtpkg.Encode(42, testSecret)
		require.NoError(t, err)

		userID, err := jwtpkg.Decode(token, testSecret)
		require.NoError(t, err)
		assert.Equal(t, uint(42), userID)
	})

	t.Run("error_invalid_token", func(t *testing.T) {
		_, err := jwtpkg.Decode("invalid.token.string", testSecret)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "認証に失敗しました。")
	})

	t.Run("error_wrong_secret", func(t *testing.T) {
		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		_, err = jwtpkg.Decode(token, "wrong-secret")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "認証に失敗しました。")
	})

	t.Run("error_expired_token", func(t *testing.T) {
		token, err := jwtpkg.EncodeWithExpiry(1, testSecret, time.Now().Add(-1*time.Hour))
		require.NoError(t, err)

		_, err = jwtpkg.Decode(token, testSecret)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "トークンの有効期限が切れています。再度ログインしてください。")
	})
}
