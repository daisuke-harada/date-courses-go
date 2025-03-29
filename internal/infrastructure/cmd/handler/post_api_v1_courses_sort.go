package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type PostApiV1CoursesSortHandler struct {}

func (h *PostApiV1CoursesSortHandler) PostApiV1CoursesSort(ctx echo.Context) error {
	// TODO: Implement your logic here
	// Example: return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}