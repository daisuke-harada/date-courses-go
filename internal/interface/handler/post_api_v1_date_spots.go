package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1DateSpotsHandler struct {
	InputPort usecase.CreateDateSpotInputPort
}

func (h *PostApiV1DateSpotsHandler) PostApiV1DateSpots(ctx echo.Context) error {
	// フォームデータを型変換（バリデーションは usecase で実行）
	genreIDStr := ctx.FormValue("genre_id")
	prefectureIDStr := ctx.FormValue("prefecture_id")
	openingTimeStr := ctx.FormValue("opening_time")
	closingTimeStr := ctx.FormValue("closing_time")

	// GenreID のパース
	genreID := 0
	if genreIDStr != "" {
		if id, err := strconv.Atoi(genreIDStr); err == nil {
			genreID = id
		}
	}

	// PrefectureID のパース
	prefectureID := 0
	if prefectureIDStr != "" {
		if id, err := strconv.Atoi(prefectureIDStr); err == nil {
			prefectureID = id
		}
	}

	// OpeningTime のパース
	var openingTime *time.Time
	if openingTimeStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, openingTimeStr); err == nil {
			openingTime = &parsedTime
		}
	}

	// ClosingTime のパース
	var closingTime *time.Time
	if closingTimeStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, closingTimeStr); err == nil {
			closingTime = &parsedTime
		}
	}

	// Image は任意フィールド
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

	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	response := openapi.DateSpotFormResponseData{
		DateSpotId: int(output.DateSpotID),
	}
	return ctx.JSON(http.StatusCreated, response)
}
