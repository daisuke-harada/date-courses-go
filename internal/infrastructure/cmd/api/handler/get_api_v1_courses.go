package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api/openapi"
	"github.com/labstack/echo/v4"
)

type GetApiV1CoursesHandler struct{}

func (h *GetApiV1CoursesHandler) GetApiV1Courses(ctx echo.Context, arg1 openapi.GetApiV1CoursesParams) error {
	// TODO: Implement your logic here
	// Example: return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
