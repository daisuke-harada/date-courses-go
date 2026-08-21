package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type DeleteCourseInputPort interface {
	Execute(context.Context, DeleteCourseInput) error
}

type DeleteCourseInput struct {
	CourseID uint
	// OperatorID は削除を実行するユーザー（トークンの currentUser）の ID です。
	OperatorID uint
}

type DeleteCourseInteractor struct {
	CourseRepository repository.CourseRepository
}

func NewDeleteCourseUsecase(courseRepository repository.CourseRepository) DeleteCourseInputPort {
	return &DeleteCourseInteractor{CourseRepository: courseRepository}
}

func (i *DeleteCourseInteractor) Execute(ctx context.Context, input DeleteCourseInput) error {
	// 削除者から見えないコース（他人の非公開コース）は存在しないものとして扱う
	course, err := i.CourseRepository.FindByID(ctx, input.CourseID, input.OperatorID)
	if err != nil {
		return apperror.NotFoundWithCause(err)
	}

	// デートコースを削除できるのは作成者だけ
	if course.UserID != input.OperatorID {
		return apperror.Forbidden("他のユーザーのデートコースは削除できません")
	}

	if err := i.CourseRepository.DeleteByID(ctx, input.CourseID); err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
