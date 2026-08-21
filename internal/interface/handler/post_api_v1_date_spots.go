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

type PostApiV1DateSpotsHandler struct {
	InputPort usecase.CreateDateSpotInputPort
}

func (h *PostApiV1DateSpotsHandler) PostApiV1DateSpots(ctx echo.Context) error {
	// デートスポットを操作できるのは管理者だけ
	if _, err := middleware.RequireAdmin(ctx); err != nil {
		return err
	}

	genreIDStr := ctx.FormValue("genre_id")
	prefectureIDStr := ctx.FormValue("prefecture_id")

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

	var imagePtr *string
	if image := ctx.FormValue("image"); image != "" {
		imagePtr = &image
	}

	input := usecase.CreateDateSpotInput{
		Name:         ctx.FormValue("name"),
		GenreID:      genreID,
		PrefectureID: prefectureID,
		CityName:     ctx.FormValue("city_name"),
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
