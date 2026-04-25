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

func TestGetUserFollowersInteractor_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}
		follower := &model.User{ID: 3, Name: "フォロワー"}
		uwr := &model.UserWithRelations{User: follower, FollowerIDs: []int{}, FollowingIDs: []int{}, Courses: []*model.Course{}, Reviews: []*model.DateSpotReview{}}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)

		relRepo := repositorymock.NewMockRelationshipRepository(ctrl)
		relRepo.EXPECT().FindFollowersByUserID(ctx, uint(1)).Return([]*model.User{follower}, nil)

		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUsersWithRelations(ctx, []*model.User{follower}).Return([]*model.UserWithRelations{uwr}, nil)

		interactor := usecase.NewGetUserFollowersUsecase(userRepo, relRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserFollowersInput{UserID: 1})

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Len(t, output.Users, 1)
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(999)).Return(nil, errors.New("not found"))

		relRepo := repositorymock.NewMockRelationshipRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		interactor := usecase.NewGetUserFollowersUsecase(userRepo, relRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserFollowersInput{UserID: 999})

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 404, statusCode)
	})

	t.Run("error_repository_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)

		relRepo := repositorymock.NewMockRelationshipRepository(ctrl)
		relRepo.EXPECT().FindFollowersByUserID(ctx, uint(1)).Return(nil, errors.New("db error"))

		userService := servicemock.NewMockUserService(ctrl)

		interactor := usecase.NewGetUserFollowersUsecase(userRepo, relRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserFollowersInput{UserID: 1})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
