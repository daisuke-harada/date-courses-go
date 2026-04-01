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
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "error"})
	}

	response := make([]openapi.AddressAndDateSpotsData, 0, len(output.DateSpots))
	for _, ds := range output.DateSpots {
		var (
			latitude  float32
			longitude float32
		)
		if ds.Latitude != nil {
			latitude = float32(*ds.Latitude)
		}
		if ds.Longitude != nil {
			longitude = float32(*ds.Longitude)
		}

		dateSpotData := openapi.DateSpotData{
			Id:        int(ds.ID),
			Name:      ds.Name,
			CreatedAt: ds.CreatedAt,
			UpdatedAt: ds.UpdatedAt,
		}
		if ds.GenreID != nil {
			dateSpotData.GenreId = *ds.GenreID
		}
		if ds.OpeningTime != nil {
			dateSpotData.OpeningTime = *ds.OpeningTime
		}
		if ds.ClosingTime != nil {
			dateSpotData.ClosingTime = *ds.ClosingTime
		}

		response = append(response, openapi.AddressAndDateSpotsData{
			Id:        int(ds.ID),
			CityName:  ds.CityName,
			Latitude:  latitude,
			Longitude: longitude,
			DateSpot:  dateSpotData,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
