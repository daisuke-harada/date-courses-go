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

func newTestUser(id uint, name string) *model.User {
	gender := model.GenderMale
	return &model.User{
		ID:     id,
		Name:   name,
		Email:  name + "@example.com",
		Gender: gender,
		Admin:  false,
	}
}

func TestCreateRelationshipInteractor_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		currentUser := newTestUser(1, "alice")
		followedUser := newTestUser(2, "bob")
		allUsers := []*model.User{currentUser, followedUser}

		currentUwr := &model.UserWithRelations{User: currentUser, FollowerIDs: []int{}, FollowingIDs: []int{2}, Courses: []*model.Course{}, Reviews: []*model.DateSpotReview{}}
		followedUwr := &model.UserWithRelations{User: followedUser, FollowerIDs: []int{1}, FollowingIDs: []int{}, Courses: []*model.Course{}, Reviews: []*model.DateSpotReview{}}
		allUwrs := []*model.UserWithRelations{currentUwr, followedUwr}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		relationshipRepo := repositorymock.NewMockRelationshipRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(currentUser, nil)
		userRepo.EXPECT().FindByID(ctx, uint(2)).Return(followedUser, nil)
		relationshipRepo.EXPECT().Create(ctx, &model.Relationship{UserID: 1, FollowID: 2}).Return(nil)
		userRepo.EXPECT().Search(ctx, (*string)(nil)).Return(allUsers, nil)
		userService.EXPECT().BuildUsersWithRelations(ctx, allUsers).Return(allUwrs, nil)
		userService.EXPECT().BuildUserWithRelations(ctx, currentUser).Return(currentUwr, nil)
		userService.EXPECT().BuildUserWithRelations(ctx, followedUser).Return(followedUwr, nil)

		interactor := usecase.NewCreateRelationshipUsecase(userRepo, relationshipRepo, userService)
		output, err := interactor.Execute(ctx, usecase.CreateRelationshipInput{
			CurrentUserID:  1,
			FollowedUserID: 2,
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, 2, len(output.Users))
		assert.Equal(t, currentUser.ID, output.CurrentUser.User.ID)
		assert.Equal(t, followedUser.ID, output.FollowedUser.User.ID)
	})

	t.Run("error_follow_self", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		relationshipRepo := repositorymock.NewMockRelationshipRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)
		// FindByID / Create は呼ばれないので EXPECT 不要

		interactor := usecase.NewCreateRelationshipUsecase(userRepo, relationshipRepo, userService)
		output, err := interactor.Execute(ctx, usecase.CreateRelationshipInput{
			CurrentUserID:  1,
			FollowedUserID: 1,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "自分自身")
	})

	t.Run("error_current_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		relationshipRepo := repositorymock.NewMockRelationshipRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		userRepo.EXPECT().FindByID(ctx, uint(999)).Return(nil, errors.New("not found"))
		// Create は呼ばれないので EXPECT 不要

		interactor := usecase.NewCreateRelationshipUsecase(userRepo, relationshipRepo, userService)
		output, err := interactor.Execute(ctx, usecase.CreateRelationshipInput{
			CurrentUserID:  999,
			FollowedUserID: 2,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("error_followed_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		currentUser := newTestUser(1, "alice")

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		relationshipRepo := repositorymock.NewMockRelationshipRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(currentUser, nil)
		userRepo.EXPECT().FindByID(ctx, uint(999)).Return(nil, errors.New("not found"))
		// Create は呼ばれないので EXPECT 不要

		interactor := usecase.NewCreateRelationshipUsecase(userRepo, relationshipRepo, userService)
		output, err := interactor.Execute(ctx, usecase.CreateRelationshipInput{
			CurrentUserID:  1,
			FollowedUserID: 999,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
