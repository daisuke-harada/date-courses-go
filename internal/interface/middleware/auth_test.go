package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
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

// currentUserHandler は認証済みなら currentUser の ID を、未認証なら "anonymous" を返します。
// 任意認証ルートで currentUser がセットされたかどうかの検証に使います。
func currentUserHandler(ctx echo.Context) error {
	user := middleware.CurrentUser(ctx)
	if user == nil {
		return ctx.String(http.StatusOK, "anonymous")
	}
	return ctx.String(http.StatusOK, strconv.FormatUint(uint64(user.ID), 10))
}

func newEchoWithAuth(t *testing.T, userRepo *repositorymock.MockUserRepository) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.HTTPErrorHandler = middleware.CustomHTTPErrorHandler
	e.Use(middleware.JWTAuthMiddleware(testSecret, userRepo))
	e.POST("/api/v1/login", currentUserHandler)
	e.POST("/api/v1/signup", dummyHandler)
	e.GET("/", dummyHandler)
	e.GET("/api/v1/top", dummyHandler)
	e.GET("/api/v1/date_spots", dummyHandler)
	e.GET("/api/v1/users", dummyHandler)
	e.GET("/api/v1/users/:id", dummyHandler)
	e.GET("/api/v1/courses", currentUserHandler)
	e.GET("/api/v1/courses/:id", currentUserHandler)
	e.GET("/api/v1/users/:user_id/followings", dummyHandler)
	e.GET("/api/v1/users/:user_id/followers", dummyHandler)
	e.POST("/api/v1/courses", dummyHandler)
	e.DELETE("/api/v1/courses/:id", dummyHandler)
	e.POST("/api/v1/date_spot_reviews", dummyHandler)
	e.PUT("/api/v1/date_spot_reviews/:id", dummyHandler)
	e.DELETE("/api/v1/date_spot_reviews/:id", dummyHandler)
	e.POST("/api/v1/date_spots", dummyHandler)
	e.PUT("/api/v1/date_spots/:id", dummyHandler)
	e.DELETE("/api/v1/date_spots/:id", dummyHandler)
	e.POST("/api/v1/relationships", dummyHandler)
	e.DELETE("/api/v1/relationships/:current_user_id/:other_user_id", dummyHandler)
	e.PUT("/api/v1/users/:id", dummyHandler)
	e.DELETE("/api/v1/users/:id", dummyHandler)
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/courses", nil)
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/courses", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("error_invalid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/courses", nil)
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/courses", nil)
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

// TestJWTAuthMiddleware_OptionalAuth は認証必須ではないルートの挙動を検証します。
// 非公開コースを作成者にだけ返すため、認証不要な GET でもトークンがあれば
// currentUser を解決する必要があります。
func TestJWTAuthMiddleware_OptionalAuth(t *testing.T) {
	t.Run("sets_current_user_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "1", rec.Body.String())
	})

	t.Run("anonymous_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "anonymous", rec.Body.String())
	})

	// 任意認証ルートでは無効なトークンを 401 にせず匿名として扱う。
	// 401 にすると、期限切れトークンを持ったままログイン画面を開いたユーザーが
	// ログイン自体を拒否され、再ログインできなくなるため。
	t.Run("anonymous_with_invalid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses", nil)
		req.Header.Set("Authorization", "Bearer invalid.token")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "anonymous", rec.Body.String())
	})

	t.Run("anonymous_with_expired_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		token, err := jwtpkg.EncodeWithExpiry(1, testSecret, time.Now().Add(-1*time.Hour))
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "anonymous", rec.Body.String())
	})

	t.Run("anonymous_when_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(nil, errors.New("not found"))

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "anonymous", rec.Body.String())
	})

	// 期限切れトークンを持ったままログインし直せることを保証する回帰テスト
	t.Run("login_is_not_rejected_with_expired_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)

		token, err := jwtpkg.EncodeWithExpiry(1, testSecret, time.Now().Add(-1*time.Hour))
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestCoursesAuthMiddleware(t *testing.T) {
	t.Run("error_post_courses_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/courses", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_post_courses_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/courses", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_delete_courses_id_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/courses/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_delete_courses_id_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/courses/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestDateSpotReviewsAuthMiddleware(t *testing.T) {
	t.Run("error_post_date_spot_reviews_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/date_spot_reviews", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_post_date_spot_reviews_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/date_spot_reviews", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_put_date_spot_reviews_id_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/date_spot_reviews/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_put_date_spot_reviews_id_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/date_spot_reviews/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_delete_date_spot_reviews_id_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/date_spot_reviews/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_delete_date_spot_reviews_id_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/date_spot_reviews/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestDateSpotsAuthMiddleware(t *testing.T) {
	t.Run("error_post_date_spots_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/date_spots", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_post_date_spots_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/date_spots", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_put_date_spots_id_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/date_spots/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_put_date_spots_id_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/date_spots/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_delete_date_spots_id_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/date_spots/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_delete_date_spots_id_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/date_spots/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestRelationshipsAuthMiddleware(t *testing.T) {
	t.Run("error_post_relationships_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/relationships", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_post_relationships_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/relationships", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_delete_relationships_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/relationships/1/2", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_delete_relationships_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/relationships/1/2", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestUsersIdAuthMiddleware(t *testing.T) {
	t.Run("error_put_users_id_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_put_users_id_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error_delete_users_id_without_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("success_delete_users_id_with_valid_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := &model.User{ID: 1, Name: "alice"}
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		token, err := jwtpkg.Encode(1, testSecret)
		require.NoError(t, err)

		e := newEchoWithAuth(t, userRepo)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
