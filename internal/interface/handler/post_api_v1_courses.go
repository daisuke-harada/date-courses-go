package handler

import (
	"net/http"
	"strconv"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1CoursesHandler struct {
	InputPort usecase.CreateCourseInputPort
}

func (h *PostApiV1CoursesHandler) PostApiV1Courses(ctx echo.Context) error {
	userID, err := strconv.Atoi(ctx.FormValue("user_id"))
	if err != nil {
		return apperror.BadRequest("user_id は整数で指定してください")
	}

	// ctx.FormValue 呼び出し時点で ParseForm 済みなので Form マップを直接参照できる
	dateSpotIDStrs := ctx.Request().Form["date_spots[]"]
	var dateSpotIDs []uint
	for _, s := range dateSpotIDStrs {
		id, err := strconv.Atoi(s)
		if err != nil {
			return apperror.BadRequest("date_spots[] は整数で指定してください")
		}
		dateSpotIDs = append(dateSpotIDs, uint(id))
	}

	travelMode := ctx.FormValue("travel_mode")
	authority := ctx.FormValue("authority")

	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.CreateCourseInput{
		UserID:      uint(userID),
		DateSpotIDs: dateSpotIDs,
		TravelMode:  travelMode,
		Authority:   authority,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, map[string]int{"course_id": int(output.CourseID)})
}
