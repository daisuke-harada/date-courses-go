package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repositorymock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	jwtpkg "github.com/daisuke-harada/date-courses-go/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const testSecret = "test-secret"

func dummyHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "ok")
}

func newEchoWithAuth(t *testing.T, userRepo *repositorymock.MockUserRepository) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.HTTPErrorHandler = middleware.CustomHTTPErrorHandler
	e.Use(middleware.JWTAuthMiddleware(testSecret, userRepo))
	e.GET("/api/v1/protected", dummyHandler)
	e.POST("/api/v1/login", dummyHandler)
	e.POST("/api/v1/signup", dummyHandler)
	e.GET("/", dummyHandler)
	e.GET("/api/v1/top", dummyHandler)
	e.GET("/api/v1/date_spots", dummyHandler)
	e.GET("/api/v1/users", dummyHandler)
	e.GET("/api/v1/users/:id", dummyHandler)
	e.GET("/api/v1/users/:userId/followings", dummyHandler)
	e.GET("/api/v1/users/:userId/followers", dummyHandler)
	return e
}

func TestJWTAuthMiddleware(t *testing.T) {
	t.Run("success_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_no_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("error_invalid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("error_expired_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		token, err := jwtpkg.EncodeWithExpiry(1, testSecret, time.Now().Add(-1*time.Hour))
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("skip_login_endpoint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("skip_signup_endpoint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("skip_public_get_endpoints", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)

		for _, path := range []string{"/", "/api/v1/top", "/api/v1/date_spots", "/api/v1/users", "/api/v1/users/1"} {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code, "path: %s should be public", path)
		}
	})

	t.Run("error_get_followings_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/followings", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_get_followings_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/followings", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_get_followers_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/followers", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_get_followers_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/followers", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
