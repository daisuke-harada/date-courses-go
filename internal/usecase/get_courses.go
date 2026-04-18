package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type GetCoursesInputPort interface {
	Execute(context.Context, GetCoursesInput) (*GetCoursesOutput, error)
}

type GetCoursesInput struct{}

type GetCoursesOutput struct {
	Courses []*model.Course
}

type GetCoursesInteractor struct {
	CourseRepository repository.CourseRepository
}

func NewGetCoursesUsecase(
	courseRepository repository.CourseRepository,
) GetCoursesInputPort {
	return &GetCoursesInteractor{
		CourseRepository: courseRepository,
	}
}

func (i *GetCoursesInteractor) Execute(ctx context.Context, input GetCoursesInput) (*GetCoursesOutput, error) {
	courses, err := i.CourseRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return &GetCoursesOutput{
		Courses: courses,
	}, nil
}
