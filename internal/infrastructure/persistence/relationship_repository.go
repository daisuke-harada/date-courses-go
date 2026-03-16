package persistence

import (
	"context"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
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

func (r *relationshipRepository) GetByID(ctx context.Context, id uint) (*model.Relationship, error) {
	var relationship model.Relationship
	if err := r.db.WithContext(ctx).First(&relationship, id).Error; err != nil {
		slog.ErrorContext(ctx, "relationshipRepository.GetByID failed", "relationship_id", id, "err", err)
		return nil, err
	}
	return &relationship, nil
}

func (r *relationshipRepository) Update(ctx context.Context, relationship *model.Relationship) error {
	if err := r.db.WithContext(ctx).Save(relationship).Error; err != nil {
		slog.ErrorContext(ctx, "relationshipRepository.Update failed", "relationship_id", relationship.ID, "err", err)
		return err
	}
	slog.InfoContext(ctx, "relationshipRepository.Update succeeded", "relationship_id", relationship.ID)
	return nil
}

func (r *relationshipRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Relationship{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "relationshipRepository.Delete failed", "relationship_id", id, "err", err)
		return err
	}
	slog.InfoContext(ctx, "relationshipRepository.Delete succeeded", "relationship_id", id)
	return nil
}
