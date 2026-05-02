package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/handler"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newTestUser(id uint, email string) *model.User {
	return &model.User{
		ID:     id,
		Name:   "テストユーザー",
		Email:  email,
		Gender: model.GenderMale,
	}
}

func TestGetApiV1CoursesHandler(t *testing.T) {
	t.Run("success_returns_200_with_courses", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		courses := []*model.Course{
			{ID: 1, TravelMode: "car", Authority: "public", User: newTestUser(1, "user1@example.com")},
			{ID: 2, TravelMode: "walk", Authority: "public", User: newTestUser(2, "user2@example.com")},
		}

		mockPort := usecasemock.NewMockGetCoursesInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetCoursesInput{}).
			Return(&usecase.GetCoursesOutput{Courses: courses}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1CoursesHandler{InputPort: mockPort}
		err := h.GetApiV1Courses(ctx, openapi.GetApiV1CoursesParams{})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, 2, len(resp))
	})

	t.Run("success_with_prefecture_id_filter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prefectureID := 13
		courses := []*model.Course{
			{ID: 1, TravelMode: "car", Authority: "public", User: newTestUser(1, "user1@example.com")},
		}

		mockPort := usecasemock.NewMockGetCoursesInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetCoursesInput{PrefectureID: &prefectureID}).
			Return(&usecase.GetCoursesOutput{Courses: courses}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses?prefecture_id=13", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1CoursesHandler{InputPort: mockPort}
		err := h.GetApiV1Courses(ctx, openapi.GetApiV1CoursesParams{PrefectureId: &prefectureID})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, 1, len(resp))
	})

	t.Run("success_returns_empty_list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetCoursesInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetCoursesInput{}).
			Return(&usecase.GetCoursesOutput{Courses: []*model.Course{}}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1CoursesHandler{InputPort: mockPort}
		err := h.GetApiV1Courses(ctx, openapi.GetApiV1CoursesParams{})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, 0, len(resp))
	})

	t.Run("error_usecase_failure_returns_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPort := usecasemock.NewMockGetCoursesInputPort(ctrl)
		mockPort.EXPECT().
			Execute(gomock.Any(), usecase.GetCoursesInput{}).
			Return(nil, apperror.InternalServerError(nil))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/courses", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		h := handler.GetApiV1CoursesHandler{InputPort: mockPort}
		err := h.GetApiV1Courses(ctx, openapi.GetApiV1CoursesParams{})

		assert.Error(t, err)
	})
}
