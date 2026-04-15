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

// FindByID は id でユーザーを検索します。
func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		slog.ErrorContext(ctx, "userRepository.FindByID failed", "err", err)
		return nil, err
	}
	return &user, nil
}

// FindByName は name でユーザーを検索します。
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

// Search は管理者を除くユーザーを名前で部分一致検索します。
func (r *userRepository) Search(ctx context.Context, name *string) ([]*model.User, error) {
	var users []*model.User
	db := r.db.WithContext(ctx).Where("admin = false")
	if name != nil && *name != "" {
		db = db.Where("name LIKE ?", "%"+*name+"%")
	}
	if err := db.Find(&users).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.Search failed", "err", err)
		return nil, err
	}
	return users, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		slog.ErrorContext(ctx, "userRepository.ExistsByEmail failed", "err", err)
		return false, err
	}
	return count > 0, nil
}

// FindFollowerIDsByUserID は指定ユーザーをフォローしているユーザーのIDリストを返します。
func (r *userRepository) FindFollowerIDsByUserID(ctx context.Context, userID uint) ([]int, error) {
	var ids []int
	if err := r.db.WithContext(ctx).
		Table("relationships").
		Where("follow_id = ?", userID).
		Pluck("user_id", &ids).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.FindFollowerIDsByUserID failed", "err", err)
		return nil, err
	}
	return ids, nil
}

// FindFollowingIDsByUserID は指定ユーザーがフォローしているユーザーのIDリストを返します。
func (r *userRepository) FindFollowingIDsByUserID(ctx context.Context, userID uint) ([]int, error) {
	var ids []int
	if err := r.db.WithContext(ctx).
		Table("relationships").
		Where("user_id = ?", userID).
		Pluck("follow_id", &ids).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.FindFollowingIDsByUserID failed", "err", err)
		return nil, err
	}
	return ids, nil
}

// Update はユーザー情報を更新します。
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.Update failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "userRepository.Update succeeded", "user_id", user.ID)
	return nil
}

// Delete はユーザーを削除します。
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.Delete failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "userRepository.Delete succeeded", "user_id", id)
	return nil
}
