package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

// GetUserInputPort はユーザー単体取得ユースケースの入力ポートです。
type GetUserInputPort interface {
	Execute(context.Context, GetUserInput) (*GetUserOutput, error)
}

type GetUserInput struct {
	ID uint
}

type GetUserOutput struct {
	UserWithRelations *model.UserWithRelations
}

type GetUserInteractor struct {
	UserRepository repository.UserRepository
	UserService    service.UserService
}

func NewGetUserUsecase(
	userRepository repository.UserRepository,
	userService service.UserService,
) GetUserInputPort {
	return &GetUserInteractor{
		UserRepository: userRepository,
		UserService:    userService,
	}
}

func (i *GetUserInteractor) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	user, err := i.UserRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	uwr, err := i.UserService.BuildUserWithRelations(ctx, user)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{UserWithRelations: uwr}, nil
}
