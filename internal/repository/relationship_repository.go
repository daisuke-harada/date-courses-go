package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type RelationshipRepository interface {
	Create(ctx context.Context, relationship *model.Relationship) error
}
