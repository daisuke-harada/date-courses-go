package usecase_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repositorymock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// デモ用アカウントの判定は名前で行うため、名前が異なれば通常のユーザーとして扱う
func TestDemoUserProtection(t *testing.T) {
	t.Run("allows_deletion_of_non_demo_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 2, Name: "alice"}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(2)).Return(user, nil)
		userRepo.EXPECT().Delete(ctx, uint(2)).Return(nil)

		interactor := usecase.NewDeleteUserUsecase(userRepo, "guest")
		require.NoError(t, interactor.Execute(ctx, usecase.DeleteUserInput{ID: 2, OperatorID: 2}))
	})

	// 設定が空なら保護しない（デモ用アカウントを置かない環境向け）
	t.Run("no_protection_when_demo_user_name_is_empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "guest"}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)
		userRepo.EXPECT().Delete(ctx, uint(1)).Return(nil)

		interactor := usecase.NewDeleteUserUsecase(userRepo, "")
		require.NoError(t, interactor.Execute(ctx, usecase.DeleteUserInput{ID: 1, OperatorID: 1}))
	})

	// 設定で別の名前を指定すれば、そのアカウントが保護される
	t.Run("protects_configured_name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 5, Name: "demo"}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(5)).Return(user, nil)

		interactor := usecase.NewDeleteUserUsecase(userRepo, "demo")
		err := interactor.Execute(ctx, usecase.DeleteUserInput{ID: 5, OperatorID: 5})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, statusCode)
	})
}
