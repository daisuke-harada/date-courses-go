package persistence

import (
	"context"
	"log/slog"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
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

// DeleteByUserIDs は user_id と follow_id の組み合わせに一致するレコードを削除します。
func (r *relationshipRepository) DeleteByUserIDs(ctx context.Context, userID uint, followID uint) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND follow_id = ?", userID, followID).
		Delete(&model.Relationship{}).Error; err != nil {
		slog.ErrorContext(ctx, "relationshipRepository.DeleteByUserIDs failed", "err", err)
		return err
	}
	return nil
}

// FindFollowingsByUserID は指定ユーザーがフォローしているユーザー一覧（管理者除く）を返します。
// Rails: user.followings.includes(...).non_admins に相当します。
func (r *relationshipRepository) FindFollowingsByUserID(ctx context.Context, userID uint) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).
		Joins("JOIN relationships ON relationships.follow_id = users.id").
		Where("relationships.user_id = ? AND users.admin = false", userID).
		Find(&users).Error; err != nil {
		slog.ErrorContext(ctx, "relationshipRepository.FindFollowingsByUserID failed", "err", err)
		return nil, err
	}
	return users, nil
}

// FindFollowersByUserID は指定ユーザーをフォローしているユーザー一覧（管理者除く）を返します。
// Rails: user.followers.includes(...).non_admins に相当します。
func (r *relationshipRepository) FindFollowersByUserID(ctx context.Context, userID uint) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).
		Joins("JOIN relationships ON relationships.user_id = users.id").
		Where("relationships.follow_id = ? AND users.admin = false", userID).
		Find(&users).Error; err != nil {
		slog.ErrorContext(ctx, "relationshipRepository.FindFollowersByUserID failed", "err", err)
		return nil, err
	}
	return users, nil
}
