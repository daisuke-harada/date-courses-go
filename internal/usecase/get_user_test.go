package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repositorymock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	servicemock "github.com/daisuke-harada/date-courses-go/internal/domain/service/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserInteractor_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}
		uwr := &model.UserWithRelations{User: user, FollowerIDs: []int{}, FollowingIDs: []int{}, Courses: []*model.Course{}, Reviews: []*model.DateSpotReview{}}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)

		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUserWithRelations(ctx, user).Return(uwr, nil)

		interactor := usecase.NewGetUserUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserInput{ID: 1})

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, uwr, output.UserWithRelations)
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(999)).Return(nil, errors.New("not found"))

		userService := servicemock.NewMockUserService(ctrl)

		interactor := usecase.NewGetUserUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserInput{ID: 999})

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 404, statusCode)
	})

	t.Run("error_service_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)

		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUserWithRelations(ctx, user).Return(nil, errors.New("service error"))

		interactor := usecase.NewGetUserUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserInput{ID: 1})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
