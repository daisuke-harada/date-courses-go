package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

// CreateRelationshipInputPort はフォロー作成ユースケースの入力ポートです。
type CreateRelationshipInputPort interface {
	Execute(context.Context, CreateRelationshipInput) (*CreateRelationshipOutput, error)
}

type CreateRelationshipInput struct {
	CurrentUserID  uint
	FollowedUserID uint
}

type CreateRelationshipOutput struct {
	Users        []*model.UserWithRelations
	CurrentUser  *model.UserWithRelations
	FollowedUser *model.UserWithRelations
}

type CreateRelationshipInteractor struct {
	UserRepository         repository.UserRepository
	RelationshipRepository repository.RelationshipRepository
	UserService            service.UserService
}

func NewCreateRelationshipUsecase(
	userRepository repository.UserRepository,
	relationshipRepository repository.RelationshipRepository,
	userService service.UserService,
) CreateRelationshipInputPort {
	return &CreateRelationshipInteractor{
		UserRepository:         userRepository,
		RelationshipRepository: relationshipRepository,
		UserService:            userService,
	}
}

func (i *CreateRelationshipInteractor) Execute(ctx context.Context, input CreateRelationshipInput) (*CreateRelationshipOutput, error) {
	// 自己フォロー禁止チェック（DBアクセス前に判定）
	if input.CurrentUserID == input.FollowedUserID {
		return nil, apperror.UnprocessableEntity("自分自身をフォローすることはできません")
	}

	// currentUser の存在確認
	currentUser, err := i.UserRepository.FindByID(ctx, input.CurrentUserID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	// followedUser の存在確認
	followedUser, err := i.UserRepository.FindByID(ctx, input.FollowedUserID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	// フォロー関係を作成
	relationship := &model.Relationship{
		UserID:   input.CurrentUserID,
		FollowID: input.FollowedUserID,
	}
	if err := i.RelationshipRepository.Create(ctx, relationship); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	// 全ユーザー一覧（non_admins）を取得
	allUsers, err := i.UserRepository.Search(ctx, nil)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	// 全ユーザーの関連データを構築
	usersWithRelations, err := i.UserService.BuildUsersWithRelations(ctx, allUsers)
	if err != nil {
		return nil, err
	}

	// currentUser / followedUser の関連データを構築
	currentUwr, err := i.UserService.BuildUserWithRelations(ctx, currentUser)
	if err != nil {
		return nil, err
	}

	followedUwr, err := i.UserService.BuildUserWithRelations(ctx, followedUser)
	if err != nil {
		return nil, err
	}

	return &CreateRelationshipOutput{
		Users:        usersWithRelations,
		CurrentUser:  currentUwr,
		FollowedUser: followedUwr,
	}, nil
}
