package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type GetCoursesInputPort interface {
	Execute(context.Context, GetCoursesInput) (*GetCoursesOutput, error)
}

type GetCoursesInput struct {
	PrefectureID *int
}

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
	courses, err := i.CourseRepository.Search(ctx, repository.CourseSearchParams{
		PrefectureID: input.PrefectureID,
	})
	if err != nil {
		return nil, err
	}

	return &GetCoursesOutput{
		Courses: courses,
	}, nil
}
