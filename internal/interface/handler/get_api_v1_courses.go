package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1CoursesHandler struct {
	InputPort usecase.GetCoursesInputPort
}

func (h *GetApiV1CoursesHandler) GetApiV1Courses(ctx echo.Context) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetCoursesInput{})
	if err != nil {
		return err
	}

	resp, err := openapi.BuildGetCoursesResponse(output)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return ctx.JSON(http.StatusOK, resp)
}
