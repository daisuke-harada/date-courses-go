package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1DateSpotsHandler struct {
	InputPort usecase.GetDateSpotsInputPort
}

// TODO: 検索ロジックの作成が必要
func (h *GetApiV1DateSpotsHandler) GetApiV1DateSpots(ctx echo.Context, params openapi.GetApiV1DateSpotsParams) error {
	input := usecase.GetDateSpotsInput{
		DateSpotName: params.DateSpotName,
		PrefectureID: params.PrefectureId,
		GenreID:      params.GenreId,
		ComeTime:     params.ComeTime,
	}
	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, openapi.NewDateSpotsResponse(output.DateSpots))
}
