package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

// UpdateUserInputPort はユーザー更新ユースケースの入力ポートです。
type UpdateUserInputPort interface {
	Execute(context.Context, UpdateUserInput) (*UpdateUserOutput, error)
}

type UpdateUserInput struct {
	ID                   uint
	Name                 string
	Email                string
	Gender               model.Gender
	Image                *string
	Password             string
	PasswordConfirmation string
}

type UpdateUserOutput struct {
	UserWithRelations *model.UserWithRelations
}

type UpdateUserInteractor struct {
	UserRepository repository.UserRepository
	UserService    service.UserService
}

func NewUpdateUserUsecase(
	userRepository repository.UserRepository,
	userService service.UserService,
) UpdateUserInputPort {
	return &UpdateUserInteractor{
		UserRepository: userRepository,
		UserService:    userService,
	}
}

func (i *UpdateUserInteractor) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	user, err := i.UserRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	// パスワードが指定されている場合のみ更新（Rails の allow_nil: true に対応）
	if err := user.ApplyUpdate(input.Name, input.Email, input.Gender, input.Image, input.Password); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if err := i.UserRepository.Update(ctx, user); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	uwr, err := i.UserService.BuildUserWithRelations(ctx, user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserOutput{UserWithRelations: uwr}, nil
}
