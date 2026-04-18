package handler

import (
	"net/http"

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
	course := output.Course
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":          course.ID,
		"user_id":     course.UserID,
		"travel_mode": course.TravelMode,
		"authority":   course.Authority,
		"created_at":  course.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":  course.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}
