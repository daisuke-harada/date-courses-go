package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

type CreateCourseInputPort interface {
	Execute(context.Context, CreateCourseInput) (*CreateCourseOutput, error)
}

type CreateCourseInput struct {
	UserID      uint
	DateSpotIDs []uint
	TravelMode  string
	Authority   string
}

func (i *CreateCourseInput) Validate() error {
	var errs []string
	if i.UserID == 0 {
		errs = append(errs, "ユーザーIDを入力してください")
	}
	if len(i.DateSpotIDs) == 0 {
		errs = append(errs, "デートスポットを1件以上入力してください")
	}
	validTravelModes := map[string]bool{"DRIVING": true, "WALKING": true}
	if !validTravelModes[i.TravelMode] {
		errs = append(errs, "移動手段はDRIVINGまたはWALKINGを指定してください")
	}
	validAuthorities := map[string]bool{"公開": true, "非公開": true}
	if !validAuthorities[i.Authority] {
		errs = append(errs, "公開設定は公開または非公開を指定してください")
	}
	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}
	return nil
}

type CreateCourseOutput struct {
	CourseID uint
}

type CreateCourseInteractor struct {
	CourseRepository     repository.CourseRepository
	DuringSpotRepository repository.DuringSpotRepository
}

func NewCreateCourseUsecase(
	courseRepo repository.CourseRepository,
	duringSpotRepo repository.DuringSpotRepository,
) CreateCourseInputPort {
	return &CreateCourseInteractor{
		CourseRepository:     courseRepo,
		DuringSpotRepository: duringSpotRepo,
	}
}

func (i *CreateCourseInteractor) Execute(ctx context.Context, input CreateCourseInput) (*CreateCourseOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	course := &model.Course{
		UserID:     input.UserID,
		TravelMode: input.TravelMode,
		Authority:  input.Authority,
	}
	if err := i.CourseRepository.Create(ctx, course); err != nil {
		return nil, apperror.InternalServerError(err)
	}
	for _, dateSpotID := range input.DateSpotIDs {
		duringSpot := &model.DuringSpot{
			CourseID:   course.ID,
			DateSpotID: dateSpotID,
		}
		if err := i.DuringSpotRepository.Create(ctx, duringSpot); err != nil {
			return nil, apperror.InternalServerError(err)
		}
	}
	return &CreateCourseOutput{CourseID: course.ID}, nil
}
