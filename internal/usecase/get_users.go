package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

// GetUsersInputPort はユーザー一覧取得ユースケースの入力ポートです。
type GetUsersInputPort interface {
	Execute(context.Context, GetUsersInput) (*GetUsersOutput, error)
}

type GetUsersInput struct {
	Name *string
}

type GetUsersOutput struct {
	Users []*model.UserWithRelations
}

type GetUsersInteractor struct {
	UserRepository repository.UserRepository
	UserService    service.UserService
}

func NewGetUsersUsecase(
	userRepository repository.UserRepository,
	userService service.UserService,
) GetUsersInputPort {
	return &GetUsersInteractor{
		UserRepository: userRepository,
		UserService:    userService,
	}
}

func (i *GetUsersInteractor) Execute(ctx context.Context, input GetUsersInput) (*GetUsersOutput, error) {
	users, err := i.UserRepository.Search(ctx, input.Name)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	result, err := i.UserService.BuildUsersWithRelations(ctx, users)
	if err != nil {
		return nil, err
	}

	return &GetUsersOutput{Users: result}, nil
}
