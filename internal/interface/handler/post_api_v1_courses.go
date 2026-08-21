package handler

import (
	"net/http"
	"strconv"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1CoursesHandler struct {
	InputPort usecase.CreateCourseInputPort
}

func (h *PostApiV1CoursesHandler) PostApiV1Courses(ctx echo.Context) error {
	// 作成者はリクエストではなくトークンから決める。
	// リクエストの user_id を信用すると他人名義のコースを登録できてしまう。
	currentUser, err := middleware.RequireCurrentUser(ctx)
	if err != nil {
		return err
	}

	form, err := ctx.FormParams()
	if err != nil {
		return apperror.BadRequestWithCause(err)
	}

	dateSpotIDStrs := form["date_spots[]"]
	var dateSpotIDs []uint
	for _, s := range dateSpotIDStrs {
		id, err := strconv.Atoi(s)
		if err != nil {
			return apperror.BadRequest("date_spots[] は整数で指定してください")
		}
		dateSpotIDs = append(dateSpotIDs, uint(id))
	}

	travelMode := form.Get("travel_mode")
	authority := form.Get("authority")

	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.CreateCourseInput{
		UserID:      currentUser.ID,
		DateSpotIDs: dateSpotIDs,
		TravelMode:  travelMode,
		Authority:   authority,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, openapi.NewCreateCourseResponse(output.CourseID))
}
