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

		courseRepo := repositorymock.NewMockCourseRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUserWithRelations(ctx, user).Return(uwr, nil)

		interactor := usecase.NewGetUserUsecase(userRepo, courseRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserInput{ID: 1})

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, uwr, output.UserWithRelations)
	})

	// 本人がマイページを開いたときだけ、非公開コース込みで取り直す
	t.Run("success_owner_gets_all_courses", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}
		publicCourses := []*model.Course{{ID: 10, UserID: 1, Authority: model.CourseAuthorityPublic}}
		allCourses := []*model.Course{
			{ID: 10, UserID: 1, Authority: model.CourseAuthorityPublic},
			{ID: 11, UserID: 1, Authority: model.CourseAuthorityPrivate},
		}
		uwr := &model.UserWithRelations{User: user, FollowerIDs: []int{}, FollowingIDs: []int{}, Courses: publicCourses, Reviews: []*model.DateSpotReview{}}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)

		courseRepo := repositorymock.NewMockCourseRepository(ctrl)
		courseRepo.EXPECT().FindAllByUserID(ctx, uint(1)).Return(allCourses, nil)

		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUserWithRelations(ctx, user).Return(uwr, nil)

		interactor := usecase.NewGetUserUsecase(userRepo, courseRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserInput{ID: 1, ViewerID: 1})

		require.NoError(t, err)
		assert.Equal(t, 2, len(output.UserWithRelations.Courses))
	})

	// 他人のマイページでは取り直さないので、公開コースだけが残る
	t.Run("success_other_viewer_gets_public_courses_only", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		user := &model.User{ID: 1, Name: "テストユーザー"}
		publicCourses := []*model.Course{{ID: 10, UserID: 1, Authority: model.CourseAuthorityPublic}}
		uwr := &model.UserWithRelations{User: user, FollowerIDs: []int{}, FollowingIDs: []int{}, Courses: publicCourses, Reviews: []*model.DateSpotReview{}}

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(1)).Return(user, nil)

		courseRepo := repositorymock.NewMockCourseRepository(ctrl)

		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUserWithRelations(ctx, user).Return(uwr, nil)

		interactor := usecase.NewGetUserUsecase(userRepo, courseRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserInput{ID: 1, ViewerID: 99})

		require.NoError(t, err)
		assert.Equal(t, 1, len(output.UserWithRelations.Courses))
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().FindByID(ctx, uint(999)).Return(nil, errors.New("not found"))

		courseRepo := repositorymock.NewMockCourseRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)

		interactor := usecase.NewGetUserUsecase(userRepo, courseRepo, userService)
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

		courseRepo := repositorymock.NewMockCourseRepository(ctrl)
		userService := servicemock.NewMockUserService(ctrl)
		userService.EXPECT().BuildUserWithRelations(ctx, user).Return(nil, errors.New("service error"))

		interactor := usecase.NewGetUserUsecase(userRepo, courseRepo, userService)
		output, err := interactor.Execute(ctx, usecase.GetUserInput{ID: 1})

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
