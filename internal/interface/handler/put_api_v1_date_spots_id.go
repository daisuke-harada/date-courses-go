package handler

import (
	"net/http"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PutApiV1DateSpotsIdHandler struct {
	InputPort usecase.UpdateDateSpotInputPort
}

func (h *PutApiV1DateSpotsIdHandler) PutApiV1DateSpotsId(ctx echo.Context, id int) error {
	var req struct {
		Name         string `form:"name" validate:"required"`
		GenreID      int    `form:"genre_id" validate:"required"`
		PrefectureID int    `form:"prefecture_id" validate:"required"`
		CityName     string `form:"city_name" validate:"required"`
		OpeningTime  *time.Time `form:"opening_time"`
		ClosingTime  *time.Time `form:"closing_time"`
		Image        *string    `form:"image"`
	}

	if err := ctx.Bind(&req); err != nil {
		return apperror.UnprocessableEntity(err)
	}

	if err := ctx.Validate(&req); err != nil {
		return apperror.UnprocessableEntity(err)
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

	if err := h.InputPort.Execute(ctx.Request().Context(), input); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
