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

func validSignupInput() usecase.SignupInput {
	return usecase.SignupInput{
		Name:                 "新規ユーザー",
		Email:                "newuser@example.com",
		Gender:               model.GenderFemale,
		Password:             "password123",
		PasswordConfirmation: "password123",
	}
}

func TestSignupInteractor_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().ExistsByEmail(ctx, "newuser@example.com").Return(false, nil)
		userRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, u *model.User) error {
			u.ID = 5
			return nil
		})

		authService := servicemock.NewMockAuthService(ctrl)
		authService.EXPECT().HashPassword("password123").Return("hashed_password", nil)

		interactor := usecase.NewSignupUsecase(userRepo, authService)
		output, err := interactor.Execute(ctx, validSignupInput())

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, uint(5), output.User.ID)
	})

	t.Run("error_validation_invalid_gender", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		authService := servicemock.NewMockAuthService(ctrl)

		input := validSignupInput()
		input.Gender = "その他" // invalid

		interactor := usecase.NewSignupUsecase(userRepo, authService)
		output, err := interactor.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 422, statusCode)
	})

	t.Run("error_email_already_exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().ExistsByEmail(ctx, "newuser@example.com").Return(true, nil)

		authService := servicemock.NewMockAuthService(ctrl)

		interactor := usecase.NewSignupUsecase(userRepo, authService)
		output, err := interactor.Execute(ctx, validSignupInput())

		assert.Error(t, err)
		assert.Nil(t, output)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, 422, statusCode)
	})

	t.Run("error_repository_create_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		userRepo := repositorymock.NewMockUserRepository(ctrl)
		userRepo.EXPECT().ExistsByEmail(ctx, "newuser@example.com").Return(false, nil)
		userRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("db error"))

		authService := servicemock.NewMockAuthService(ctrl)
		authService.EXPECT().HashPassword("password123").Return("hashed_password", nil)

		interactor := usecase.NewSignupUsecase(userRepo, authService)
		output, err := interactor.Execute(ctx, validSignupInput())

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
