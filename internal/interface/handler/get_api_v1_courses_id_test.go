package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetApiV1CoursesIdHandler(t *testing.T) {
	t.Run("success_returns_200_with_course", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetCourseInput{CourseID: 1}).
			Return(&usecase.GetCourseOutput{
				Course: &model.Course{
					ID:         1,
					UserID:     2,
					TravelMode: "DRIVING",
					Authority:  "公開",
				},
			}, nil)

		ctx, rec := setupGetRequest("/api/v1/courses/1")

		h := handler.GetApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.GetApiV1CoursesId(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, float64(1), resp["id"])
		assert.Equal(t, float64(2), resp["user_id"])
		assert.Equal(t, "DRIVING", resp["travel_mode"])
		assert.Equal(t, "公開", resp["authority"])
	})

	t.Run("error_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetCourseInput{CourseID: 999}).
			Return(nil, apperror.NotFound())

		ctx, _ := setupGetRequest("/api/v1/courses/999")

		h := handler.GetApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.GetApiV1CoursesId(ctx, 999)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("error_usecase_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetCourseInput{CourseID: 1}).
			Return(nil, apperror.InternalServerError(errors.New("db error")))

		ctx, _ := setupGetRequest("/api/v1/courses/1")

		h := handler.GetApiV1CoursesIdHandler{InputPort: mockPort}
		err := h.GetApiV1CoursesId(ctx, 1)

		assert.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
