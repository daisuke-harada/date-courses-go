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

func validUpdateUserInput() usecase.UpdateUserInput {
	return usecase.UpdateUserInput{
		ID:     1,
		Name:   "更新ユーザー",
		Email:  "updated@example.com",
		Gender: model.GenderMale,
	}
}

func TestUpdateUserInteractor_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "元ユーザー", Email: "old@example.com", Gender: model.GenderFemale}
		uwr := &model.UserWithRelations{User: user, FollowerIDs: []int{}, FollowingIDs: []int{}, Courses: []*model.Course{}, Reviews: []*model.DateSpotReview{}}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)
		userRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUserWithRelations(ctx, gomock.Any()).Return(uwr, nil)

		interactor := usecase.NewUpdateUserUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, validUpdateUserInput())

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, uwr, output.UserWithRelations)
	})

	t.Run("error_validation_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		input := validUpdateUserInput()
		input.Name = "" // invalid

		interactor := usecase.NewUpdateUserUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 422, statusCode)
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(nil, errors.New("not found"))

		userService := servicemock.NewMockUserService(ctrl)

		interactor := usecase.NewUpdateUserUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, validUpdateUserInput())

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 404, statusCode)
	})

	t.Run("error_repository_update_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "元ユーザー", Email: "old@example.com", Gender: model.GenderFemale}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)
		userRepo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("db error"))

		userService := servicemock.NewMockUserService(ctrl)

		interactor := usecase.NewUpdateUserUsecase(userRepo, userService)
		output, err := interactor.Execute(ctx, validUpdateUserInput())

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 500, statusCode)
	})
}
