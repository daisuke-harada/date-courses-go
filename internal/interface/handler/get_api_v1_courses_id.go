package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1CoursesIdHandler struct {
	InputPort usecase.GetCourseInputPort
}

func (h *GetApiV1CoursesIdHandler) GetApiV1CoursesId(ctx echo.Context, arg1 int) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetCourseInput{CourseID: uint(arg1)})
	if err != nil {
		return err
	}

	resp, err := openapi.NewCourseResponse(output.Course)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return ctx.JSON(http.StatusOK, resp)
}
