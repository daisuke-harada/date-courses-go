package persistence

import (
	"context"
	"log/slog"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
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

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.GetByID failed", "user_id", id, "err", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.Update failed", "user_id", user.ID, "err", err)
		return err
	}
	slog.InfoContext(ctx, "userRepository.Update succeeded", "user_id", user.ID)
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.Delete failed", "user_id", id, "err", err)
		return err
	}
	slog.InfoContext(ctx, "userRepository.Delete succeeded", "user_id", id)
	return nil
}
