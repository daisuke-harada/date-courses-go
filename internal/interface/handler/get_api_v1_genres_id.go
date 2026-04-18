package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1GenresIdHandler struct {
	InputPort usecase.GetDateSpotsInputPort
}

func (h *GetApiV1GenresIdHandler) GetApiV1GenresId(ctx echo.Context, arg1 int) error {
	genreID := arg1
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetDateSpotsInput{
		GenreID: &genreID,
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"address_and_date_spots": openapi.NewDateSpotsResponse(output.DateSpots),
	})
}
