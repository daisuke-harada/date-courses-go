package openapi

import "github.com/daisuke-harada/date-courses-go/internal/domain/model"

// SignupResponseBody は POST /api/v1/signup の成功レスポンスです。
// Rails の RegistrationSerializer と同じ構造を再現します。
type SignupResponseBody struct {
	User        UserResponseData `json:"user"`
	LoginStatus bool             `json:"login_status"`
}

// NewSignupResponse は model.User から SignupResponseBody を生成します。
// 新規ユーザーのため followerIds, followingIds, courses, date_spot_reviews は空配列です。
func NewSignupResponse(user *model.User) SignupResponseBody {
	return SignupResponseBody{
		User: UserResponseData{
			Id:              int(user.ID),
			Admin:           user.Admin,
			Email:           user.Email,
			Gender:          user.Gender,
			Image:           ImageData{Url: user.Image},
			Name:            user.Name,
			FollowerIds:     []int{},
			FollowingIds:    []int{},
			Courses:         []CourseResponseData{},
			DateSpotReviews: []DateSpotReviewData{},
		},
		LoginStatus: true,
	}
}
