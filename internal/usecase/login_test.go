package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repositorymock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	servicemock "github.com/daisuke-harada/date-courses-go/internal/domain/service/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

const testJWTSecret usecase.JWTSecretKey = "test-secret"

func newLoginUser() *model.User {
	gender := model.GenderMale
	return &model.User{
		ID:             1,
		Name:           "alice",
		Email:          "alice@example.com",
		Gender:         gender,
		PasswordDigest: "$2a$10$hashed",
	}
}

func TestLoginInteractor_Execute(t *testing.T) {
	t.Run("success_returns_user_and_token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := newLoginUser()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		authService := servicemock.NewMockAuthService(ctrl)

		userRepo.EXPECT().FindByName(ctx, "alice").Return(user, nil)
		authService.EXPECT().CheckPassword(user.PasswordDigest, "password").Return(true)

		interactor := usecase.NewLoginUsecase(userRepo, authService, testJWTSecret)
		output, err := interactor.Execute(ctx, usecase.LoginInput{Name: "alice", Password: "password"})

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, user, output.User)
		assert.NotEmpty(t, output.Token)
	})

	t.Run("error_empty_name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		authService := servicemock.NewMockAuthService(ctrl)

		interactor := usecase.NewLoginUsecase(userRepo, authService, testJWTSecret)
		output, err := interactor.Execute(context.Background(), usecase.LoginInput{Name: "", Password: "password"})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "名前を入力してください")
	})

	t.Run("error_empty_password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		authService := servicemock.NewMockAuthService(ctrl)

		interactor := usecase.NewLoginUsecase(userRepo, authService, testJWTSecret)
		output, err := interactor.Execute(context.Background(), usecase.LoginInput{Name: "alice", Password: ""})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "パスワードを入力してください")
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		authService := servicemock.NewMockAuthService(ctrl)

		userRepo.EXPECT().FindByName(ctx, "unknown").Return(nil, gorm.ErrRecordNotFound)

		interactor := usecase.NewLoginUsecase(userRepo, authService, testJWTSecret)
		output, err := interactor.Execute(ctx, usecase.LoginInput{Name: "unknown", Password: "password"})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "認証に失敗しました。")
	})

	t.Run("error_wrong_password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := newLoginUser()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		authService := servicemock.NewMockAuthService(ctrl)

		userRepo.EXPECT().FindByName(ctx, "alice").Return(user, nil)
		authService.EXPECT().CheckPassword(user.PasswordDigest, "wrong").Return(false)

		interactor := usecase.NewLoginUsecase(userRepo, authService, testJWTSecret)
		output, err := interactor.Execute(ctx, usecase.LoginInput{Name: "alice", Password: "wrong"})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "認証に失敗しました。")
	})

	t.Run("error_db_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		userRepo := repositorymock.NewMockUserRepository(ctrl)
		authService := servicemock.NewMockAuthService(ctrl)

		userRepo.EXPECT().FindByName(ctx, "alice").Return(nil, errors.New("db error"))

		interactor := usecase.NewLoginUsecase(userRepo, authService, testJWTSecret)
		output, err := interactor.Execute(ctx, usecase.LoginInput{Name: "alice", Password: "password"})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
