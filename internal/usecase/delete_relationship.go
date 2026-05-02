package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

type DeleteRelationshipInputPort interface {
	Execute(context.Context, DeleteRelationshipInput) (*DeleteRelationshipOutput, error)
}

type DeleteRelationshipInput struct {
	UserID   uint
	FollowID uint
}

type DeleteRelationshipOutput struct {
	Users          []*model.UserWithRelations
	CurrentUser    *model.UserWithRelations
	UnfollowedUser *model.UserWithRelations
}

type DeleteRelationshipInteractor struct {
	RelationshipRepository repository.RelationshipRepository
	UserRepository         repository.UserRepository
	UserService            service.UserService
}

func NewDeleteRelationshipUsecase(
	relationshipRepository repository.RelationshipRepository,
	userRepository repository.UserRepository,
	userService service.UserService,
) DeleteRelationshipInputPort {
	return &DeleteRelationshipInteractor{
		RelationshipRepository: relationshipRepository,
		UserRepository:         userRepository,
		UserService:            userService,
	}
}

func (i *DeleteRelationshipInteractor) Execute(ctx context.Context, input DeleteRelationshipInput) (*DeleteRelationshipOutput, error) {
	currentUser, err := i.UserRepository.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	unfollowedUser, err := i.UserRepository.FindByID(ctx, input.FollowID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	if err := i.RelationshipRepository.DeleteByUserIDs(ctx, input.UserID, input.FollowID); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	allUsers, err := i.UserRepository.Search(ctx, nil)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	usersWithRelations, err := i.UserService.BuildUsersWithRelations(ctx, allUsers)
	if err != nil {
		return nil, err
	}

	currentUwr, err := i.UserService.BuildUserWithRelations(ctx, currentUser)
	if err != nil {
		return nil, err
	}

	unfollowedUwr, err := i.UserService.BuildUserWithRelations(ctx, unfollowedUser)
	if err != nil {
		return nil, err
	}

	return &DeleteRelationshipOutput{
		Users:          usersWithRelations,
		CurrentUser:    currentUwr,
		UnfollowedUser: unfollowedUwr,
	}, nil
}
