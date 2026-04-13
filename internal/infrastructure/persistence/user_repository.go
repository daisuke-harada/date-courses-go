package persistence

import (
	"context"
	"errors"
	"log/slog"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "userRepository.Create succeeded", "user_id", user.ID)
	return nil
}

// FindByName は name でユーザーを検索します。
// 見つからない場合は gorm.ErrRecordNotFound を返します。
func (r *userRepository) FindByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		slog.ErrorContext(ctx, "userRepository.FindByName failed", "err", err)
		return nil, err
	}
	slog.InfoContext(ctx, "userRepository.FindByName succeeded", "user_id", user.ID)
	return &user, nil
}
