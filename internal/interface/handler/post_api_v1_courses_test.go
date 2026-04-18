package handler_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPostApiV1CoursesHandler(t *testing.T) {
	t.Run("success_returns_201_with_course_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.CreateCourseInput{
				UserID:      1,
				DateSpotIDs: []uint{10, 20},
				TravelMode:  "DRIVING",
				Authority:   "公開",
			}).
			Return(&usecase.CreateCourseOutput{CourseID: 5}, nil)

		form := url.Values{}
		form.Set("user_id", "1")
		form.Add("date_spots[]", "10")
		form.Add("date_spots[]", "20")
		form.Set("travel_mode", "DRIVING")
		form.Set("authority", "公開")
		ctx, rec := setupFormRequest(http.MethodPost, "/api/v1/courses", form)

		h := handler.PostApiV1CoursesHandler{InputPort: mockPort}
		err := h.PostApiV1Courses(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, float64(5), resp["course_id"])
	})

	t.Run("error_bad_request_invalid_user_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateCourseInputPort(ctrl)

		form := url.Values{}
		form.Set("user_id", "abc")
		form.Add("date_spots[]", "10")
		form.Set("travel_mode", "DRIVING")
		form.Set("authority", "公開")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/courses", form)

		h := handler.PostApiV1CoursesHandler{InputPort: mockPort}
		err := h.PostApiV1Courses(ctx)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("error_bad_request_invalid_date_spot_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateCourseInputPort(ctrl)

		form := url.Values{}
		form.Set("user_id", "1")
		form.Add("date_spots[]", "abc")
		form.Set("travel_mode", "DRIVING")
		form.Set("authority", "公開")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/courses", form)

		h := handler.PostApiV1CoursesHandler{InputPort: mockPort}
		err := h.PostApiV1Courses(ctx)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("error_usecase_returns_unprocessable_entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.UnprocessableEntity("デートスポットを1件以上入力してください"))

		form := url.Values{}
		form.Set("user_id", "1")
		form.Set("travel_mode", "DRIVING")
		form.Set("authority", "公開")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/courses", form)

		h := handler.PostApiV1CoursesHandler{InputPort: mockPort}
		err := h.PostApiV1Courses(ctx)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("error_usecase_returns_internal_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockCreateCourseInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, apperror.InternalServerError(nil))

		form := url.Values{}
		form.Set("user_id", "1")
		form.Add("date_spots[]", "10")
		form.Set("travel_mode", "DRIVING")
		form.Set("authority", "公開")
		ctx, _ := setupFormRequest(http.MethodPost, "/api/v1/courses", form)

		h := handler.PostApiV1CoursesHandler{InputPort: mockPort}
		err := h.PostApiV1Courses(ctx)

		require.Error(t, err)
		statusCode, _, _, ok := apperror.HTTPStatus(err)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
