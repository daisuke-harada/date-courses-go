package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type DeleteApiV1CoursesIdHandler struct {
	InputPort usecase.DeleteCourseInputPort
}

func (h *DeleteApiV1CoursesIdHandler) DeleteApiV1CoursesId(ctx echo.Context, courseID int) error {
	// 削除できるのは作成者だけなので、操作主体はトークンから決める
	currentUser, err := middleware.RequireCurrentUser(ctx)
	if err != nil {
		return err
	}

	if err := h.InputPort.Execute(ctx.Request().Context(), usecase.DeleteCourseInput{
		CourseID:   uint(courseID),
		OperatorID: currentUser.ID,
	}); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
