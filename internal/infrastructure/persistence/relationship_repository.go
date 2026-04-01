package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"gorm.io/gorm"
)

type relationshipRepository struct {
	db *gorm.DB
}

func NewRelationshipRepository(db *gorm.DB) repository.RelationshipRepository {
	return &relationshipRepository{db: db}
}

func (r *relationshipRepository) Create(ctx context.Context, relationship *model.Relationship) error {
	if err := r.db.WithContext(ctx).Create(relationship).Error; err != nil {
		slog.ErrorContext(ctx, "relationshipRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "relationshipRepository.Create succeeded", "relationship_id", relationship.ID)
	return nil
}
