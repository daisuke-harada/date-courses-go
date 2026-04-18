package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type DeleteCourseInputPort interface {
	Execute(context.Context, DeleteCourseInput) (*DeleteCourseOutput, error)
}

type DeleteCourseInput struct {
	CourseID uint
}

type DeleteCourseOutput struct{}

type DeleteCourseInteractor struct {
	CourseRepository repository.CourseRepository
}

func NewDeleteCourseUsecase(courseRepository repository.CourseRepository) DeleteCourseInputPort {
	return &DeleteCourseInteractor{CourseRepository: courseRepository}
}

func (i *DeleteCourseInteractor) Execute(ctx context.Context, input DeleteCourseInput) (*DeleteCourseOutput, error) {
	if err := i.CourseRepository.DeleteByID(ctx, input.CourseID); err != nil {
		return nil, apperror.InternalServerError(err)
	}
	return &DeleteCourseOutput{}, nil
}
