package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain"
)

type RelationshipRepository interface {
	Create(ctx context.Context, relationship *domain.Relationship) error
	GetByID(ctx context.Context, id uint) (*domain.Relationship, error)
	Update(ctx context.Context, relationship *domain.Relationship) error
	Delete(ctx context.Context, id uint) error
}
