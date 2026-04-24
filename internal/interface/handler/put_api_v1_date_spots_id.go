package handler

import (
	"net/http"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PutApiV1DateSpotsIdHandler struct {
	InputPort usecase.UpdateDateSpotInputPort
}

func (h *PutApiV1DateSpotsIdHandler) PutApiV1DateSpotsId(ctx echo.Context, id int) error {
	var req struct {
		Name         string     `form:"name"`
		GenreID      int        `form:"genre_id"`
		PrefectureID int        `form:"prefecture_id"`
		CityName     string     `form:"city_name"`
		OpeningTime  *time.Time `form:"opening_time"`
		ClosingTime  *time.Time `form:"closing_time"`
		Image        *string    `form:"image"`
	}

	if err := ctx.Bind(&req); err != nil {
		// 型変換エラーはバインド失敗として処理
		return err
	}

	input := usecase.UpdateDateSpotInput{
		DateSpotID:   uint(id),
		Name:         req.Name,
		GenreID:      req.GenreID,
		PrefectureID: req.PrefectureID,
		CityName:     req.CityName,
		OpeningTime:  req.OpeningTime,
		ClosingTime:  req.ClosingTime,
		Image:        req.Image,
	}

	if err := input.Validate(); err != nil {
		return err
	}

	if err := h.InputPort.Execute(ctx.Request().Context(), input); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
