package usecase_test

import (
	"context"
	"errors"
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

func TestDeleteUserInteractor_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)
		userRepo.EXPECT().Delete(ctx, uint(1)).Return(nil)

		interactor := usecase.NewDeleteUserUsecase(userRepo)
		err := interactor.Execute(ctx, usecase.DeleteUserInput{ID: 1, OperatorID: 1})

		require.NoError(t, err)
	})

	// 他人のアカウントを削除できないことを保証する
	t.Run("error_forbidden_when_operator_is_not_the_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)

		interactor := usecase.NewDeleteUserUsecase(userRepo)
		err := interactor.Execute(ctx, usecase.DeleteUserInput{ID: 1, OperatorID: 99})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, statusCode)
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(999)).Return(nil, errors.New("not found"))

		interactor := usecase.NewDeleteUserUsecase(userRepo)
		err := interactor.Execute(ctx, usecase.DeleteUserInput{ID: 999, OperatorID: 999})

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 404, statusCode)
	})

	t.Run("error_repository_delete_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)
		userRepo.EXPECT().Delete(ctx, uint(1)).Return(errors.New("db error"))

		interactor := usecase.NewDeleteUserUsecase(userRepo)
		err := interactor.Execute(ctx, usecase.DeleteUserInput{ID: 1, OperatorID: 1})

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 500, statusCode)
	})
}
