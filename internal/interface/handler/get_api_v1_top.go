package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1TopHandler struct {
	InputPort usecase.GetDateSpotsInputPort
}

func (h *GetApiV1TopHandler) GetApiV1Top(ctx echo.Context) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetDateSpotsInput{})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, openapi.NewTopResponse(output.DateSpots))
}
