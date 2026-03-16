package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type RelationshipRepository interface {
	Create(ctx context.Context, relationship *model.Relationship) error
	GetByID(ctx context.Context, id uint) (*model.Relationship, error)
	Update(ctx context.Context, relationship *model.Relationship) error
	Delete(ctx context.Context, id uint) error
}
