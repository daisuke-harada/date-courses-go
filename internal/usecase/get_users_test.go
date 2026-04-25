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
)

func TestGetUsersInteractor_Execute(t *testing.T) {
	t.Run("success_returns_users", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		users := []*model.User{
			{ID: 1, Name: "ユーザー1"},
			{ID: 2, Name: "ユーザー2"},
		}
		uwrs := []*model.UserWithRelations{
			{User: users[0], FollowerIDs: []int{}, FollowingIDs: []int{}, Courses: []*model.Course{}, Reviews: []*model.DateSpotReview{}},
			{User: users[1], FollowerIDs: []int{}, FollowingIDs: []int{}, Courses: []*model.Course{}, Reviews: []*model.DateSpotReview{}},
		}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().Search(ctx, (*string)(nil)).Return(users, nil)

		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUsersWithRelations(ctx, users).Return(uwrs, nil)

		interactor := usecase.NewGetUsersUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUsersInput{})

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Len(t, output.Users, 2)
	})

	t.Run("error_repository_search_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().Search(ctx, (*string)(nil)).Return(nil, errors.New("db error"))

		userService := servicemock.NewMockUserService(ctrl)

		interactor := usecase.NewGetUsersUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUsersInput{})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
