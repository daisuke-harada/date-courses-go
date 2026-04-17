package openapi

import (
	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/oapi-codegen/runtime/types"
)

// UserResponseBody は UserResponseData の代替型です。
// CourseResponseDataBody / DateSpotReviewDataBody を使うことで
// OpeningTime / ClosingTime の nullable を正しく扱います。
type UserResponseBody struct {
	Admin           bool                     `json:"admin"`
	Courses         []CourseResponseDataBody `json:"courses"`
	DateSpotReviews []DateSpotReviewDataBody `json:"date_spot_reviews"`
	Email           string                   `json:"email"`
	FollowerIds     []int                    `json:"followerIds"`
	FollowingIds    []int                    `json:"followingIds"`
	Gender          Gender                   `json:"gender"`
	Id              int                      `json:"id"`
	Image           ImageData                `json:"image"`
	Name            string                   `json:"name"`
}

// CourseResponseDataBody は CourseResponseData の代替型です。
// DateSpots フィールドに AddressAndDateSpotsDataResponse を使います。
type CourseResponseDataBody struct {
	Authority                  string                            `json:"authority"`
	DateSpots                  []AddressAndDateSpotsDataResponse `json:"date_spots"`
	Id                         int                               `json:"id"`
	NoDuplicatePrefectureNames []string                          `json:"no_duplicate_prefecture_names"`
	TravelMode                 string                            `json:"travel_mode"`
	User                       UserData                          `json:"user"`
}

// DateSpotReviewDataBody は DateSpotReviewData の代替型です。
// DateSpot フィールドに DateSpotDataResponse を使います。
type DateSpotReviewDataBody struct {
	Content  string               `json:"content"`
	DateSpot DateSpotDataResponse `json:"date_spot"`
	Id       int                  `json:"id"`
	Rate     float32              `json:"rate"`
}

// BuildUserResponseBody は User と関連データから UserResponseBody を構築します。
func BuildUserResponseBody(
	user *model.User,
	followerIDs, followingIDs []int,
	courses []*model.Course,
	reviews []*model.DateSpotReview,
) (UserResponseBody, error) {
	gender, err := NewGender(user.Gender)
	if err != nil {
		return UserResponseBody{}, err
	}

	courseResponses := make([]CourseResponseDataBody, 0, len(courses))
	for _, c := range courses {
		cr, err := buildCourseResponseBody(c)
		if err != nil {
			return UserResponseBody{}, err
		}
		courseResponses = append(courseResponses, cr)
	}

	reviewResponses := make([]DateSpotReviewDataBody, 0, len(reviews))
	for _, rv := range reviews {
		reviewResponses = append(reviewResponses, buildDateSpotReviewDataBody(rv))
	}

	if followerIDs == nil {
		followerIDs = []int{}
	}
	if followingIDs == nil {
		followingIDs = []int{}
	}

	return UserResponseBody{
		Id:              int(user.ID),
		Admin:           user.Admin,
		Email:           user.Email,
		Gender:          gender,
		Image:           ImageData{Url: user.Image},
		Name:            user.Name,
		FollowerIds:     followerIDs,
		FollowingIds:    followingIDs,
		Courses:         courseResponses,
		DateSpotReviews: reviewResponses,
	}, nil
}

// buildCourseResponseBody は Course モデルから CourseResponseDataBody を構築します。
// Course.User と Course.DuringSpots.DateSpot が Preload 済みであることを前提とします。
func buildCourseResponseBody(course *model.Course) (CourseResponseDataBody, error) {
	dateSpots := make([]AddressAndDateSpotsDataResponse, 0, len(course.DuringSpots))
	prefectureIDSet := make(map[int]struct{})

	for _, ds := range course.DuringSpots {
		if ds.DateSpot == nil {
			continue
		}
		dateSpots = append(dateSpots, newAddressAndDateSpotsData(ds.DateSpot))
		if ds.DateSpot.PrefectureID != nil {
			prefectureIDSet[*ds.DateSpot.PrefectureID] = struct{}{}
		}
	}

	prefectureNames := make([]string, 0, len(prefectureIDSet))
	for id := range prefectureIDSet {
		name := master.PrefectureNameByID(id)
		if name != "" {
			prefectureNames = append(prefectureNames, name)
		}
	}

	var courseUser UserData
	if course.User != nil {
		gender, err := NewGender(course.User.Gender)
		if err != nil {
			return CourseResponseDataBody{}, err
		}
		courseUser = UserData{
			Id:     int(course.User.ID),
			Name:   course.User.Name,
			Email:  types.Email(course.User.Email),
			Gender: gender,
			Image:  ImageData{Url: course.User.Image},
			Admin:  course.User.Admin,
		}
	}

	return CourseResponseDataBody{
		Id:                         int(course.ID),
		Authority:                  course.Authority,
		TravelMode:                 course.TravelMode,
		DateSpots:                  dateSpots,
		NoDuplicatePrefectureNames: prefectureNames,
		User:                       courseUser,
	}, nil
}

// buildDateSpotReviewDataBody は DateSpotReview から DateSpotReviewDataBody を構築します。
// Review.DateSpot が Preload 済みであることを前提とします。
func buildDateSpotReviewDataBody(review *model.DateSpotReview) DateSpotReviewDataBody {
	var rate float32
	if review.Rate != nil {
		rate = float32(*review.Rate)
	}
	var content string
	if review.Content != nil {
		content = *review.Content
	}

	var dateSpot DateSpotDataResponse
	if review.DateSpot != nil {
		dateSpot = newDateSpotData(review.DateSpot)
	}

	return DateSpotReviewDataBody{
		Id:       int(review.ID),
		Rate:     rate,
		Content:  content,
		DateSpot: dateSpot,
	}
}

// NewGetUsersResponse は output.Users から []UserResponseBody を構築します。
func NewGetUsersResponse(users []*model.UserWithRelations) ([]UserResponseBody, error) {
	responses := make([]UserResponseBody, 0, len(users))
	for _, uwr := range users {
		resp, err := BuildUserResponseBody(
			uwr.User,
			uwr.FollowerIDs,
			uwr.FollowingIDs,
			uwr.Courses,
			uwr.Reviews,
		)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

// NewUserWithRelationsResponse は *model.UserWithRelations から UserResponseBody を構築します。
func NewUserWithRelationsResponse(uwr *model.UserWithRelations) (UserResponseBody, error) {
	return BuildUserResponseBody(uwr.User, uwr.FollowerIDs, uwr.FollowingIDs, uwr.Courses, uwr.Reviews)
}
