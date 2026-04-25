package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type DeleteRelationshipInputPort interface {
	Execute(context.Context, DeleteRelationshipInput) error
}

type DeleteRelationshipInput struct {
	UserID   uint
	FollowID uint
}

type DeleteRelationshipInteractor struct {
	RelationshipRepository repository.RelationshipRepository
}

func NewDeleteRelationshipUsecase(relationshipRepository repository.RelationshipRepository) DeleteRelationshipInputPort {
	return &DeleteRelationshipInteractor{RelationshipRepository: relationshipRepository}
}

func (i *DeleteRelationshipInteractor) Execute(ctx context.Context, input DeleteRelationshipInput) error {
	if err := i.RelationshipRepository.DeleteByUserIDs(ctx, input.UserID, input.FollowID); err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
