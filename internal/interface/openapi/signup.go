package openapi

import (
	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

// SignupResponseBody は POST /api/v1/signup の成功レスポンスです。
// Rails の RegistrationSerializer と同じ構造を再現します。
type SignupResponseBody struct {
	User        UserResponseData `json:"user"`
	LoginStatus bool             `json:"login_status"`
}

// NewSignupResponse は model.User から SignupResponseBody を生成します。
// 新規ユーザーのため followerIds, followingIds, courses, date_spot_reviews は空配列です。
func NewSignupResponse(user *model.User) (SignupResponseBody, error) {
	gender, err := NewGender(user.Gender)
	if err != nil {
		return SignupResponseBody{}, err
	}
	return SignupResponseBody{
		User: UserResponseData{
			Id:              int(user.ID),
			Admin:           user.Admin,
			Email:           user.Email,
			Gender:          gender,
			Image:           ImageData{Url: user.Image},
			Name:            user.Name,
			FollowerIds:     []int{},
			FollowingIds:    []int{},
			Courses:         []CourseResponseData{},
			DateSpotReviews: []DateSpotReviewData{},
		},
		LoginStatus: true,
	}, nil
}

// NewGender は model.Gender を openapi.Gender に変換します。
// 不正な値の場合は error を返します。
func NewGender(g model.Gender) (Gender, error) {
	switch g {
	case model.GenderMale:
		return GenderEmpty, nil
	case model.GenderFemale:
		return GenderN1, nil
	default:
		return "", apperror.InternalServerError(nil)
	}
}
