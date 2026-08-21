package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// DeleteUserInputPort はユーザー削除ユースケースの入力ポートです。
type DeleteUserInputPort interface {
	Execute(context.Context, DeleteUserInput) error
}

type DeleteUserInput struct {
	ID uint
	// OperatorID は削除を実行するユーザー（トークンの currentUser）の ID です。
	OperatorID uint
}

type DeleteUserInteractor struct {
	UserRepository repository.UserRepository
}

func NewDeleteUserUsecase(userRepository repository.UserRepository) DeleteUserInputPort {
	return &DeleteUserInteractor{
		UserRepository: userRepository,
	}
}

func (i *DeleteUserInteractor) Execute(ctx context.Context, input DeleteUserInput) error {
	user, err := i.UserRepository.FindByID(ctx, input.ID)
	if err != nil {
		return apperror.NotFound()
	}

	// 退会できるのは本人だけ
	if user.ID != input.OperatorID {
		return apperror.Forbidden("他のユーザーのアカウントは削除できません")
	}

	if err := i.UserRepository.Delete(ctx, user.ID); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}
