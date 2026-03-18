package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1DateSpotsHandler struct {
	InputPort usecase.GetDateSpotsInputPort
}

func (h *GetApiV1DateSpotsHandler) GetApiV1DateSpots(ctx echo.Context) error {
	output, err := h.InputPort.Execute(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "error"})
	}

	response := make([]openapi.AddressAndDateSpotsData, 0, len(output.DateSpots))
	for i, ds := range output.DateSpots {
		addr := output.Addresses[i]

		var (
			cityName  string
			latitude  float32
			longitude float32
		)
		if addr != nil {
			cityName = addr.CityName
			if addr.Latitude != nil {
				latitude = float32(*addr.Latitude)
			}
			if addr.Longitude != nil {
				longitude = float32(*addr.Longitude)
			}
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
			Id: func() int {
				if addr != nil {
					return int(addr.ID)
				}
				return 0
			}(),
			CityName:  cityName,
			Latitude:  latitude,
			Longitude: longitude,
			DateSpot:  dateSpotData,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
