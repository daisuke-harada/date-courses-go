package openapi

import (
	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
)

// NewUserResponseData は User と関連データから UserResponseBody を構築します。
func NewUserResponseData(
	user *model.User,
	followerIDs, followingIDs []int,
	courses []*model.Course,
	reviews []*model.DateSpotReview,
) (UserResponseData, error) {
	gender, err := NewGender(user.Gender)
	if err != nil {
		return UserResponseData{}, err
	}

	courseResponses := make([]CourseResponseData, 0, len(courses))
	for _, c := range courses {
		cr, err := buildCourseResponseBody(c)
		if err != nil {
			return UserResponseData{}, err
		}
		courseResponses = append(courseResponses, cr)
	}

	reviewResponses := lo.Map(reviews, func(rv *model.DateSpotReview, _ int) DateSpotReviewData {
		return newDateSpotReviewData(rv)
	})

	if followerIDs == nil {
		followerIDs = []int{}
	}
	if followingIDs == nil {
		followingIDs = []int{}
	}

	return UserResponseData{
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

func buildCourseResponseBody(course *model.Course) (CourseResponseData, error) {
	dateSpots := make([]DateSpotSummaryData, 0, len(course.DuringSpots))
	prefectureIDSet := make(map[int]struct{})

	for _, ds := range course.DuringSpots {
		if ds.DateSpot == nil {
			continue
		}
		dateSpots = append(dateSpots, newDateSpotSummaryData(ds.DateSpot))
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
			return CourseResponseData{}, err
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

	return CourseResponseData{
		Id:                         int(course.ID),
		Authority:                  course.Authority,
		TravelMode:                 course.TravelMode,
		DateSpots:                  dateSpots,
		NoDuplicatePrefectureNames: prefectureNames,
		User:                       courseUser,
	}, nil
}

func newDateSpotReviewData(review *model.DateSpotReview) DateSpotReviewData {
	var rate float32
	if review.Rate != nil {
		rate = float32(*review.Rate)
	}
	var content string
	if review.Content != nil {
		content = *review.Content
	}

	var dateSpot DateSpotData
	if review.DateSpot != nil {
		dateSpot = newDateSpotData(review.DateSpot)
	}

	return DateSpotReviewData{
		Id:       int(review.ID),
		Rate:     rate,
		Content:  content,
		DateSpot: dateSpot,
	}
}

func NewGetUsersResponse(users []*model.UserWithRelations) ([]UserResponseData, error) {
	responses := make([]UserResponseData, 0, len(users))
	for _, uwr := range users {
		resp, err := NewUserResponseData(
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
func NewUserWithRelationsResponse(uwr *model.UserWithRelations) (UserResponseData, error) {
	return NewUserResponseData(uwr.User, uwr.FollowerIDs, uwr.FollowingIDs, uwr.Courses, uwr.Reviews)
}
