package handler

import (
	"net/http"
	"strconv"
	"strings"
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
	var errs []string

	name := ctx.FormValue("name")
	genreIDStr := ctx.FormValue("genre_id")
	prefectureIDStr := ctx.FormValue("prefecture_id")
	cityName := ctx.FormValue("city_name")
	openingTimeStr := ctx.FormValue("opening_time")
	closingTimeStr := ctx.FormValue("closing_time")
	image := ctx.FormValue("image")

	// name: presence
	if strings.TrimSpace(name) == "" {
		errs = append(errs, "スポット名を入力してください")
	}

	// genre_id: presence, integer
	var genreID int
	if genreIDStr == "" {
		errs = append(errs, "ジャンルを選択してください")
	} else {
		id, err := strconv.Atoi(genreIDStr)
		if err != nil {
			errs = append(errs, "ジャンルは整数値で指定してください")
		} else {
			genreID = id
		}
	}

	// prefecture_id: presence, integer
	var prefectureID int
	if prefectureIDStr == "" {
		errs = append(errs, "都道府県を選択してください")
	} else {
		id, err := strconv.Atoi(prefectureIDStr)
		if err != nil {
			errs = append(errs, "都道府県は整数値で指定してください")
		} else {
			prefectureID = id
		}
	}

	// city_name: presence
	if strings.TrimSpace(cityName) == "" {
		errs = append(errs, "市区町村を入力してください")
	}

	// opening_time: presence, RFC3339 format
	var openingTime *time.Time
	if openingTimeStr == "" {
		errs = append(errs, "営業開始時刻を入力してください")
	} else {
		if parsedTime, err := time.Parse(time.RFC3339, openingTimeStr); err != nil {
			errs = append(errs, "営業開始時刻の形式が正しくありません")
		} else {
			openingTime = &parsedTime
		}
	}

	// closing_time: presence, RFC3339 format
	var closingTime *time.Time
	if closingTimeStr == "" {
		errs = append(errs, "営業終了時刻を入力してください")
	} else {
		if parsedTime, err := time.Parse(time.RFC3339, closingTimeStr); err != nil {
			errs = append(errs, "営業終了時刻の形式が正しくありません")
		} else {
			closingTime = &parsedTime
		}
	}

	// image: optional
	var imagePtr *string
	if image != "" {
		imagePtr = &image
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}

	input := usecase.CreateDateSpotInput{
		Name:         name,
		GenreID:      genreID,
		PrefectureID: prefectureID,
		CityName:     cityName,
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
