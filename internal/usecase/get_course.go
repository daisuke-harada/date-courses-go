package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type GetCourseInputPort interface {
	Execute(context.Context, GetCourseInput) (*GetCourseOutput, error)
}

type GetCourseInput struct {
	CourseID uint
}

type GetCourseOutput struct {
	Course *model.Course
}

type GetCourseInteractor struct {
	CourseRepository repository.CourseRepository
}

func NewGetCourseUsecase(courseRepository repository.CourseRepository) GetCourseInputPort {
	return &GetCourseInteractor{
		CourseRepository: courseRepository,
	}
}

func (i *GetCourseInteractor) Execute(ctx context.Context, input GetCourseInput) (*GetCourseOutput, error) {
	course, err := i.CourseRepository.FindByID(ctx, input.CourseID)
	if err != nil {
		return nil, apperror.NotFound()
	}
	return &GetCourseOutput{Course: course}, nil
}
