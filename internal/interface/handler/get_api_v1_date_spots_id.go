package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1DateSpotsIdHandler struct {
	InputPort usecase.GetDateSpotInputPort
}

func (h *GetApiV1DateSpotsIdHandler) GetApiV1DateSpotsId(ctx echo.Context, id int) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetDateSpotInput{
		ID: uint(id),
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, openapi.NewDateSpotResponse(output.DateSpot))
}
