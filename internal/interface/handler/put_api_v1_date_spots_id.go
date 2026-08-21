package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PutApiV1DateSpotsIdHandler struct {
	InputPort usecase.UpdateDateSpotInputPort
}

func (h *PutApiV1DateSpotsIdHandler) PutApiV1DateSpotsId(ctx echo.Context, id int) error {
	// デートスポットを操作できるのは管理者だけ
	if _, err := middleware.RequireAdmin(ctx); err != nil {
		return err
	}

	var req struct {
		Name         string  `form:"name"`
		GenreID      int     `form:"genre_id"`
		PrefectureID int     `form:"prefecture_id"`
		CityName     string  `form:"city_name"`
		Image        *string `form:"image"`
	}

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	input := usecase.UpdateDateSpotInput{
		DateSpotID:   uint(id),
		Name:         req.Name,
		GenreID:      req.GenreID,
		PrefectureID: req.PrefectureID,
		CityName:     req.CityName,
		Image:        req.Image,
	}

	if err := h.InputPort.Execute(ctx.Request().Context(), input); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, openapi.DateSpotFormResponseData{DateSpotId: id})
}
