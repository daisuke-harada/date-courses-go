package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1PrefecturesIdHandler struct {
	InputPort usecase.GetDateSpotsInputPort
}

func (h *GetApiV1PrefecturesIdHandler) GetApiV1PrefecturesId(ctx echo.Context, arg1 int) error {
	prefectureID := arg1
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetDateSpotsInput{
		PrefectureID: &prefectureID,
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, openapi.NewDateSpotsResponse(output.DateSpots))
}
