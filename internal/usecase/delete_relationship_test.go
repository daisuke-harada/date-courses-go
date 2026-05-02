package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repomock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	servicemock "github.com/daisuke-harada/date-courses-go/internal/domain/service/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteRelationshipInteractor_Execute(t *testing.T) {
	ctx := context.Background()

	currentUser := &model.User{ID: 1, Name: "current", Email: "current@example.com", Gender: "男性"}
	unfollowedUser := &model.User{ID: 2, Name: "unfollowed", Email: "unfollowed@example.com", Gender: "女性"}
	uwr := &model.UserWithRelations{User: currentUser, FollowerIDs: []int{}, FollowingIDs: []int{}}
	unfollowedUwr := &model.UserWithRelations{User: unfollowedUser, FollowerIDs: []int{}, FollowingIDs: []int{}}

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		relationshipRepo := repomock.NewMockRelationshipRepository(ctrl)
		userRepo := repomock.NewMockUserRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(currentUser, nil)
		userRepo.EXPECT().FindByID(ctx, uint(2)).Return(unfollowedUser, nil)
		relationshipRepo.EXPECT().DeleteByUserIDs(ctx, uint(1), uint(2)).Return(nil)
		userRepo.EXPECT().Search(ctx, nil).Return([]*model.User{currentUser, unfollowedUser}, nil)
		userService.EXPECT().
			BuildUsersWithRelations(ctx, gomock.Any()).
			Return([]*model.UserWithRelations{uwr, unfollowedUwr}, nil)
		userService.EXPECT().BuildUserWithRelations(ctx, currentUser).Return(uwr, nil)
		userService.EXPECT().BuildUserWithRelations(ctx, unfollowedUser).Return(unfollowedUwr, nil)

		interactor := usecase.NewDeleteRelationshipUsecase(relationshipRepo, userRepo, userService)
		output, err := interactor.Execute(ctx, usecase.DeleteRelationshipInput{
			UserID:   1,
			FollowID: 2,
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotNil(t, output.CurrentUser)
		assert.NotNil(t, output.UnfollowedUser)
	})

	t.Run("error_current_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		relationshipRepo := repomock.NewMockRelationshipRepository(ctrl)
		userRepo := repomock.NewMockUserRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(nil, errors.New("not found"))

		interactor := usecase.NewDeleteRelationshipUsecase(relationshipRepo, userRepo, userService)
		output, err := interactor.Execute(ctx, usecase.DeleteRelationshipInput{
			UserID:   1,
			FollowID: 2,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("error_repository_delete_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		relationshipRepo := repomock.NewMockRelationshipRepository(ctrl)
		userRepo := repomock.NewMockUserRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(currentUser, nil)
		userRepo.EXPECT().FindByID(ctx, uint(2)).Return(unfollowedUser, nil)
		relationshipRepo.EXPECT().DeleteByUserIDs(ctx, uint(1), uint(2)).Return(errors.New("db error"))

		interactor := usecase.NewDeleteRelationshipUsecase(relationshipRepo, userRepo, userService)
		output, err := interactor.Execute(ctx, usecase.DeleteRelationshipInput{
			UserID:   1,
			FollowID: 2,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
