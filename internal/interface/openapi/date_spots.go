package openapi

import (
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/samber/lo"
)

// NewDateSpotShowResponse は DateSpot とレビュー一覧から DateSpotShowResponse を構築します。
func NewDateSpotShowResponse(dateSpot *model.DateSpot, reviews []*model.DateSpotReview) DateSpotShowResponseData {
	return DateSpotShowResponseData{
		DateSpot:          newDateSpotSummaryData(dateSpot),
		ReviewAverageRate: float32(dateSpot.AverageRate),
		DateSpotReviews:   newDateSpotShowResponseDataDateSpotReviewsInner(reviews),
	}
}

func newDateSpotShowResponseDataDateSpotReviewsInner(reviews []*model.DateSpotReview) []DateSpotShowResponseDataDateSpotReviewsInner {
	return lo.Map(reviews, func(r *model.DateSpotReview, _ int) DateSpotShowResponseDataDateSpotReviewsInner {
		item := DateSpotShowResponseDataDateSpotReviewsInner{
			Id:         int(r.ID),
			DateSpotId: int(r.DateSpotID),
			UserId:     int(r.UserID),
		}
		if r.Rate != nil {
			f := float32(*r.Rate)
			item.Rate = &f
		}
		if r.Content != nil {
			item.Content = r.Content
		}
		if r.User != nil {
			item.UserName = r.User.Name
			item.UserGender = string(r.User.Gender)
			item.UserImage = ImageData{Url: r.User.Image}
		}
		return item
	})
}

func NewCreateDateSpotResponse(dateSpotID uint) DateSpotFormResponseData {
	return DateSpotFormResponseData{
		DateSpotId: int(dateSpotID),
	}
}

func NewDateSpotResponse(dateSpot *model.DateSpot) DateSpotSummaryData {
	return newDateSpotSummaryData(dateSpot)
}

func NewDateSpotsResponse(dateSpots []*model.DateSpot) []DateSpotSummaryData {
	return lo.Map(dateSpots, func(ds *model.DateSpot, _ int) DateSpotSummaryData {
		return newDateSpotSummaryData(ds)
	})
}

func NewDateSpotSummaries(dateSpots []*model.DateSpot) []DateSpotSummaryData {
	return lo.Map(dateSpots, func(ds *model.DateSpot, _ int) DateSpotSummaryData {
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

		return DateSpotSummaryData{
			AverageRate:       float32(ds.AverageRate),
			CityName:          ds.CityName,
			DateSpot:          dateSpot,
			GenreName:         genreName,
			Id:                int(ds.ID),
			Latitude:          latitude,
			Longitude:         longitude,
			PrefectureName:    prefectureName,
			ReviewTotalNumber: ds.ReviewTotalNumber,
		}
	})
}

func newDateSpotSummaryData(ds *model.DateSpot) DateSpotSummaryData {
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

	return DateSpotSummaryData{
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

func newDateSpotData(ds *model.DateSpot) DateSpotData {
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

	return DateSpotData{
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
}
