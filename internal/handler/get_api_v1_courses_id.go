package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func (h *handler) GetApiV1CoursesId(ctx echo.Context, arg1 int ) error {
	// TODO: Implement your logic here
	// Example: return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}