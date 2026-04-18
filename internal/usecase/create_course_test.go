package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	repomock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateCourseUsecase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockDuringSpotRepo := repomock.NewMockDuringSpotRepository(ctrl)

		mockCourseRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, course *model.Course) error {
				course.ID = 10
				return nil
			})
		mockDuringSpotRepo.EXPECT().
			Create(gomock.Any(), &model.DuringSpot{CourseID: 10, DateSpotID: 1}).
			Return(nil)
		mockDuringSpotRepo.EXPECT().
			Create(gomock.Any(), &model.DuringSpot{CourseID: 10, DateSpotID: 2}).
			Return(nil)

		uc := usecase.NewCreateCourseUsecase(mockCourseRepo, mockDuringSpotRepo)
		out, err := uc.Execute(context.Background(), usecase.CreateCourseInput{
			UserID:      1,
			DateSpotIDs: []uint{1, 2},
			TravelMode:  "DRIVING",
			Authority:   "公開",
		})

		require.NoError(t, err)
		assert.Equal(t, uint(10), out.CourseID)
	})

	t.Run("error_validation_missing_user_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockDuringSpotRepo := repomock.NewMockDuringSpotRepository(ctrl)

		uc := usecase.NewCreateCourseUsecase(mockCourseRepo, mockDuringSpotRepo)
		_, err := uc.Execute(context.Background(), usecase.CreateCourseInput{
			UserID:      0,
			DateSpotIDs: []uint{1},
			TravelMode:  "DRIVING",
			Authority:   "公開",
		})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_validation_no_date_spots", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockDuringSpotRepo := repomock.NewMockDuringSpotRepository(ctrl)

		uc := usecase.NewCreateCourseUsecase(mockCourseRepo, mockDuringSpotRepo)
		_, err := uc.Execute(context.Background(), usecase.CreateCourseInput{
			UserID:      1,
			DateSpotIDs: []uint{},
			TravelMode:  "DRIVING",
			Authority:   "公開",
		})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_validation_invalid_travel_mode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockDuringSpotRepo := repomock.NewMockDuringSpotRepository(ctrl)

		uc := usecase.NewCreateCourseUsecase(mockCourseRepo, mockDuringSpotRepo)
		_, err := uc.Execute(context.Background(), usecase.CreateCourseInput{
			UserID:      1,
			DateSpotIDs: []uint{1},
			TravelMode:  "FLYING",
			Authority:   "公開",
		})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_validation_invalid_authority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockDuringSpotRepo := repomock.NewMockDuringSpotRepository(ctrl)

		uc := usecase.NewCreateCourseUsecase(mockCourseRepo, mockDuringSpotRepo)
		_, err := uc.Execute(context.Background(), usecase.CreateCourseInput{
			UserID:      1,
			DateSpotIDs: []uint{1},
			TravelMode:  "DRIVING",
			Authority:   "未設定",
		})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_course_repository_create_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockDuringSpotRepo := repomock.NewMockDuringSpotRepository(ctrl)

		mockCourseRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		uc := usecase.NewCreateCourseUsecase(mockCourseRepo, mockDuringSpotRepo)
		_, err := uc.Execute(context.Background(), usecase.CreateCourseInput{
			UserID:      1,
			DateSpotIDs: []uint{1},
			TravelMode:  "DRIVING",
			Authority:   "公開",
		})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("error_during_spot_repository_create_failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCourseRepo := repomock.NewMockCourseRepository(ctrl)
		mockDuringSpotRepo := repomock.NewMockDuringSpotRepository(ctrl)

		mockCourseRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, course *model.Course) error {
				course.ID = 10
				return nil
			})
		mockDuringSpotRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		uc := usecase.NewCreateCourseUsecase(mockCourseRepo, mockDuringSpotRepo)
		_, err := uc.Execute(context.Background(), usecase.CreateCourseInput{
			UserID:      1,
			DateSpotIDs: []uint{1},
			TravelMode:  "DRIVING",
			Authority:   "公開",
		})

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
