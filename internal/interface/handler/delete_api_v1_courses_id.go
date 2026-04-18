package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type DeleteApiV1CoursesIdHandler struct {
	InputPort usecase.DeleteCourseInputPort
}

func (h *DeleteApiV1CoursesIdHandler) DeleteApiV1CoursesId(ctx echo.Context, arg1 int) error {
	_, err := h.InputPort.Execute(ctx.Request().Context(), usecase.DeleteCourseInput{
		CourseID: uint(arg1),
	})
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
