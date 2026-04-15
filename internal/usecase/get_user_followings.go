package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

// GetUserFollowingsInputPort はユーザーのフォロー一覧取得ユースケースの入力ポートです。
type GetUserFollowingsInputPort interface {
	Execute(context.Context, GetUserFollowingsInput) (*GetUserFollowingsOutput, error)
}

type GetUserFollowingsInput struct {
	UserID uint
}

type GetUserFollowingsOutput struct {
	Users []*model.UserWithRelations
}

type GetUserFollowingsInteractor struct {
	UserRepository         repository.UserRepository
	RelationshipRepository repository.RelationshipRepository
	UserService            service.UserService
}

func NewGetUserFollowingsUsecase(
	userRepository repository.UserRepository,
	relationshipRepository repository.RelationshipRepository,
	userService service.UserService,
) GetUserFollowingsInputPort {
	return &GetUserFollowingsInteractor{
		UserRepository:         userRepository,
		RelationshipRepository: relationshipRepository,
		UserService:            userService,
	}
}

func (i *GetUserFollowingsInteractor) Execute(ctx context.Context, input GetUserFollowingsInput) (*GetUserFollowingsOutput, error) {
	// ユーザーの存在確認
	user, err := i.UserRepository.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	followings, err := i.RelationshipRepository.FindFollowingsByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	result, err := i.UserService.BuildUsersWithRelations(ctx, followings)
	if err != nil {
		return nil, err
	}

	return &GetUserFollowingsOutput{Users: result}, nil
}
