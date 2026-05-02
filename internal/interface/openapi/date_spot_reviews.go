package openapi

import "github.com/daisuke-harada/date-courses-go/internal/domain/model"

// NewDateSpotReviewResponse はレビュー一覧から DateSpotReviewResponseData を構築します。
func NewDateSpotReviewResponse(reviews []*model.DateSpotReview) DateSpotReviewResponseData {
	return DateSpotReviewResponseData{
		DateSpotReviews:   toDateSpotReviewInners(reviews),
		ReviewAverageRate: computeReviewAverageRate(reviews),
	}
}

func toDateSpotReviewInners(reviews []*model.DateSpotReview) []DateSpotShowResponseDataDateSpotReviewsInner {
	items := make([]DateSpotShowResponseDataDateSpotReviewsInner, 0, len(reviews))
	for _, r := range reviews {
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
		items = append(items, item)
	}
	return items
}

func computeReviewAverageRate(reviews []*model.DateSpotReview) float32 {
	var sum float64
	var count int
	for _, r := range reviews {
		if r.Rate != nil {
			sum += *r.Rate
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return float32(sum / float64(count))
}
