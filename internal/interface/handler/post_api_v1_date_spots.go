package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1DateSpotsHandler struct {
	InputPort usecase.CreateDateSpotInputPort
}

func (h *PostApiV1DateSpotsHandler) PostApiV1DateSpots(ctx echo.Context) error {
	genreIDStr := ctx.FormValue("genre_id")
	prefectureIDStr := ctx.FormValue("prefecture_id")
	openingTimeStr := ctx.FormValue("opening_time")
	closingTimeStr := ctx.FormValue("closing_time")

	genreID := 0
	if genreIDStr != "" {
		id, err := strconv.Atoi(genreIDStr)
		if err != nil {
			return apperror.BadRequest("genre_id は数値で指定してください")
		}
		genreID = id
	}

	prefectureID := 0
	if prefectureIDStr != "" {
		id, err := strconv.Atoi(prefectureIDStr)
		if err != nil {
			return apperror.BadRequest("prefecture_id は数値で指定してください")
		}
		prefectureID = id
	}

	var openingTime *time.Time
	if openingTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, openingTimeStr)
		if err != nil {
			return apperror.UnprocessableEntity("opening_time の形式が正しくありません（RFC3339 形式で指定してください）")
		}
		openingTime = &parsedTime
	}

	var closingTime *time.Time
	if closingTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, closingTimeStr)
		if err != nil {
			return apperror.UnprocessableEntity("closing_time の形式が正しくありません（RFC3339 形式で指定してください）")
		}
		closingTime = &parsedTime
	}

	var imagePtr *string
	if image := ctx.FormValue("image"); image != "" {
		imagePtr = &image
	}

	input := usecase.CreateDateSpotInput{
		Name:         ctx.FormValue("name"),
		GenreID:      genreID,
		PrefectureID: prefectureID,
		CityName:     ctx.FormValue("city_name"),
		OpeningTime:  openingTime,
		ClosingTime:  closingTime,
		Image:        imagePtr,
	}

	if err := input.Validate(); err != nil {
		return err
	}

	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, openapi.NewCreateDateSpotResponse(output.DateSpotID))
}
