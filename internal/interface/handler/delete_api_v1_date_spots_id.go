package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type DeleteApiV1DateSpotsIdHandler struct {
	InputPort usecase.DeleteDateSpotInputPort
}

func (h *DeleteApiV1DateSpotsIdHandler) DeleteApiV1DateSpotsId(ctx echo.Context, id int) error {
	input := usecase.DeleteDateSpotInput{DateSpotID: uint(id)}
	if err := h.InputPort.Execute(ctx.Request().Context(), input); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
