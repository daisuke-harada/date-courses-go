package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type RelationshipRepository interface {
	Create(ctx context.Context, relationship *model.Relationship) error
	FindFollowingsByUserID(ctx context.Context, userID uint) ([]*model.User, error)
	FindFollowersByUserID(ctx context.Context, userID uint) ([]*model.User, error)
}
