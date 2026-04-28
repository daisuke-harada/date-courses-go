package openapi

import (
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

// DateSpotDataResponse は生成型 DateSpotData の代替で、
// OpeningTime / ClosingTime を *time.Time にして null をそのまま JSON 出力します。
type DateSpotDataResponse struct {
	AverageRate float32    `json:"average_rate"`
	ClosingTime *time.Time `json:"closing_time"`
	CreatedAt   time.Time  `json:"created_at"`
	GenreId     int        `json:"genre_id"`
	Id          int        `json:"id"`
	Image       ImageData  `json:"image"`
	Name        string     `json:"name"`
	OpeningTime *time.Time `json:"opening_time"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DateSpotSummaryDataResponse は生成型 DateSpotSummaryData の代替で、
// DateSpot フィールドに DateSpotDataResponse を使います。
type DateSpotSummaryDataResponse struct {
	AverageRate       float32              `json:"average_rate"`
	CityName          string               `json:"city_name"`
	DateSpot          DateSpotDataResponse `json:"date_spot"`
	GenreName         string               `json:"genre_name"`
	Id                int                  `json:"id"`
	Latitude          float32              `json:"latitude"`
	Longitude         float32              `json:"longitude"`
	PrefectureName    string               `json:"prefecture_name"`
	ReviewTotalNumber int                  `json:"review_total_number"`
}

// NewCreateDateSpotResponse は DateSpotID から DateSpotFormResponseData を構築します。
func NewCreateDateSpotResponse(dateSpotID uint) DateSpotFormResponseData {
	return DateSpotFormResponseData{
		DateSpotId: int(dateSpotID),
	}
}

func NewDateSpotResponse(dateSpot *model.DateSpot) DateSpotSummaryDataResponse {
	return newDateSpotSummaryData(dateSpot)
}

func NewDateSpotsResponse(dateSpots []*model.DateSpot) []DateSpotSummaryDataResponse {
	response := make([]DateSpotSummaryDataResponse, 0, len(dateSpots))
	for _, ds := range dateSpots {
		response = append(response, newDateSpotSummaryData(ds))
	}
	return response
}

// NewDateSpotSummaries は生成型の []DateSpotSummaryData を返すヘルパーです。
// Top のレスポンス（generated types）と互換にするために使います。
func NewDateSpotSummaries(dateSpots []*model.DateSpot) []DateSpotSummaryData {
	response := make([]DateSpotSummaryData, 0, len(dateSpots))
	for _, ds := range dateSpots {
		var (
			latitude       float32
			longitude      float32
			genreName      string
			prefectureName string
		)
		if ds.Latitude != nil {
			latitude = float32(*ds.Latitude)
		}
		if ds.Longitude != nil {
			longitude = float32(*ds.Longitude)
		}
		if ds.GenreID != nil {
			genreName = master.GenreNameByID(*ds.GenreID)
		}
		if ds.PrefectureID != nil {
			prefectureName = master.PrefectureNameByID(*ds.PrefectureID)
		}

		var genreId int
		if ds.GenreID != nil {
			genreId = *ds.GenreID
		}

		var openingTime, closingTime time.Time
		if ds.OpeningTime != nil {
			openingTime = *ds.OpeningTime
		}
		if ds.ClosingTime != nil {
			closingTime = *ds.ClosingTime
		}

		dateSpot := DateSpotData{
			Id:          int(ds.ID),
			Name:        ds.Name,
			Image:       ImageData{Url: ds.Image},
			GenreId:     genreId,
			AverageRate: float32(ds.AverageRate),
			CreatedAt:   ds.CreatedAt,
			UpdatedAt:   ds.UpdatedAt,
			OpeningTime: openingTime,
			ClosingTime: closingTime,
		}

		response = append(response, DateSpotSummaryData{
			AverageRate:       float32(ds.AverageRate),
			CityName:          ds.CityName,
			DateSpot:          dateSpot,
			GenreName:         genreName,
			Id:                int(ds.ID),
			Latitude:          latitude,
			Longitude:         longitude,
			PrefectureName:    prefectureName,
			ReviewTotalNumber: ds.ReviewTotalNumber,
		})
	}
	return response
}

func newDateSpotSummaryData(ds *model.DateSpot) DateSpotSummaryDataResponse {
	var (
		latitude       float32
		longitude      float32
		genreName      string
		prefectureName string
	)
	if ds.Latitude != nil {
		latitude = float32(*ds.Latitude)
	}
	if ds.Longitude != nil {
		longitude = float32(*ds.Longitude)
	}
	if ds.GenreID != nil {
		genreName = master.GenreNameByID(*ds.GenreID)
	}
	if ds.PrefectureID != nil {
		prefectureName = master.PrefectureNameByID(*ds.PrefectureID)
	}

	return DateSpotSummaryDataResponse{
		Id:                int(ds.ID),
		CityName:          ds.CityName,
		Latitude:          latitude,
		Longitude:         longitude,
		GenreName:         genreName,
		PrefectureName:    prefectureName,
		AverageRate:       float32(ds.AverageRate),
		ReviewTotalNumber: ds.ReviewTotalNumber,
		DateSpot:          newDateSpotData(ds),
	}
}

func newDateSpotData(ds *model.DateSpot) DateSpotDataResponse {
	var genreId int
	if ds.GenreID != nil {
		genreId = *ds.GenreID
	}

	return DateSpotDataResponse{
		Id:          int(ds.ID),
		Name:        ds.Name,
		Image:       ImageData{Url: ds.Image},
		GenreId:     genreId,
		AverageRate: float32(ds.AverageRate),
		CreatedAt:   ds.CreatedAt,
		UpdatedAt:   ds.UpdatedAt,
		OpeningTime: ds.OpeningTime,
		ClosingTime: ds.ClosingTime,
	}
}
