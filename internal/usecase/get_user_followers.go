package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

// GetUserFollowersInputPort はユーザーのフォロワー一覧取得ユースケースの入力ポートです。
type GetUserFollowersInputPort interface {
	Execute(context.Context, GetUserFollowersInput) (*GetUserFollowersOutput, error)
}

type GetUserFollowersInput struct {
	UserID uint
}

type GetUserFollowersOutput struct {
	Users []*model.UserWithRelations
}

type GetUserFollowersInteractor struct {
	UserRepository         repository.UserRepository
	RelationshipRepository repository.RelationshipRepository
	UserService            service.UserService
}

func NewGetUserFollowersUsecase(
	userRepository repository.UserRepository,
	relationshipRepository repository.RelationshipRepository,
	userService service.UserService,
) GetUserFollowersInputPort {
	return &GetUserFollowersInteractor{
		UserRepository:         userRepository,
		RelationshipRepository: relationshipRepository,
		UserService:            userService,
	}
}

func (i *GetUserFollowersInteractor) Execute(ctx context.Context, input GetUserFollowersInput) (*GetUserFollowersOutput, error) {
	user, err := i.UserRepository.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	followers, err := i.RelationshipRepository.FindFollowersByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	uwr, err := i.UserService.BuildUsersWithRelations(ctx, followers)
	if err != nil {
		return nil, err
	}

	return &GetUserFollowersOutput{Users: uwr}, nil
}
