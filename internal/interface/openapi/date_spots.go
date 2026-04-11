package openapi

import (
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

// DateSpotFlatResponse は AddressSerializer フラット化後のレスポンス型です。
// date_spot ネストを廃止し、全フィールドをトップレベルに配置します。
type DateSpotFlatResponse struct {
	Id                int        `json:"id"`
	Name              string     `json:"name"`
	GenreId           int        `json:"genre_id"`
	Image             ImageData  `json:"image"`
	OpeningTime       *time.Time `json:"opening_time"`
	ClosingTime       *time.Time `json:"closing_time"`
	AverageRate       float32    `json:"average_rate"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CityName          string     `json:"city_name"`
	Latitude          float32    `json:"latitude"`
	Longitude         float32    `json:"longitude"`
	PrefectureName    string     `json:"prefecture_name"`
	GenreName         string     `json:"genre_name"`
	ReviewTotalNumber int        `json:"review_total_number"`
}

func NewDateSpotsResponse(dateSpots []*model.DateSpot) []DateSpotFlatResponse {
	response := make([]DateSpotFlatResponse, 0, len(dateSpots))
	for _, ds := range dateSpots {
		response = append(response, newDateSpotFlatResponse(ds))
	}
	return response
}

func newDateSpotFlatResponse(ds *model.DateSpot) DateSpotFlatResponse {
	var (
		genreId        int
		latitude       float32
		longitude      float32
		genreName      string
		prefectureName string
	)
	if ds.GenreID != nil {
		genreId = *ds.GenreID
		genreName = master.GenreNameByID(*ds.GenreID)
	}
	if ds.Latitude != nil {
		latitude = float32(*ds.Latitude)
	}
	if ds.Longitude != nil {
		longitude = float32(*ds.Longitude)
	}
	if ds.PrefectureID != nil {
		prefectureName = master.PrefectureNameByID(*ds.PrefectureID)
	}

	return DateSpotFlatResponse{
		Id:                int(ds.ID),
		Name:              ds.Name,
		GenreId:           genreId,
		Image:             ImageData{Url: ds.Image},
		OpeningTime:       ds.OpeningTime,
		ClosingTime:       ds.ClosingTime,
		AverageRate:       float32(ds.AverageRate),
		CreatedAt:         ds.CreatedAt,
		UpdatedAt:         ds.UpdatedAt,
		CityName:          ds.CityName,
		Latitude:          latitude,
		Longitude:         longitude,
		PrefectureName:    prefectureName,
		GenreName:         genreName,
		ReviewTotalNumber: ds.ReviewTotalNumber,
	}
}
