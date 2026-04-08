package openapi

import (
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

func NewDateSpotsResponse(dateSpots []*model.DateSpot) []AddressAndDateSpotsData {
	response := make([]AddressAndDateSpotsData, 0, len(dateSpots))
	for _, ds := range dateSpots {
		response = append(response, newAddressAndDateSpotsData(ds))
	}
	return response
}

func newAddressAndDateSpotsData(ds *model.DateSpot) AddressAndDateSpotsData {
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

	dateSpotData := newDateSpotData(ds)

	return AddressAndDateSpotsData{
		Id:        int(ds.ID),
		CityName:  ds.CityName,
		Latitude:  latitude,
		Longitude: longitude,
		DateSpot:  dateSpotData,
	}
}

func newDateSpotData(ds *model.DateSpot) DateSpotData {
	var genreId int
	if ds.GenreID != nil {
		genreId = *ds.GenreID
	}

	result := DateSpotData{
		Id:      int(ds.ID),
		Name:    ds.Name,
		Image:   ImageData{Url: ds.Image},
		GenreId: genreId,
		// AverageRate は repository 側で JOIN して取得する想定のためゼロ値
		CreatedAt: ds.CreatedAt,
		UpdatedAt: ds.UpdatedAt,
	}

	if ds.OpeningTime != nil {
		result.OpeningTime = *ds.OpeningTime
	}
	if ds.ClosingTime != nil {
		result.ClosingTime = *ds.ClosingTime
	}

	return result
}
