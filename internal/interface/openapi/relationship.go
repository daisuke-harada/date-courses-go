package openapi

import "github.com/daisuke-harada/date-courses-go/internal/usecase"

func NewUnFollowResponseData(output *usecase.DeleteRelationshipOutput) (UnFollowResponseData, error) {
	users := make([]UserResponseData, 0, len(output.Users))
	for _, uwr := range output.Users {
		resp, err := NewUserResponseData(uwr.User, uwr.FollowerIDs, uwr.FollowingIDs, uwr.Courses, uwr.Reviews)
		if err != nil {
			return UnFollowResponseData{}, err
		}
		users = append(users, resp)
	}

	currentUser, err := NewUserWithRelationsResponse(output.CurrentUser)
	if err != nil {
		return UnFollowResponseData{}, err
	}

	unfollowedUser, err := NewUserWithRelationsResponse(output.UnfollowedUser)
	if err != nil {
		return UnFollowResponseData{}, err
	}

	return UnFollowResponseData{
		Users:          users,
		CurrentUser:    currentUser,
		UnfollowedUser: unfollowedUser,
	}, nil
}

func NewFollowResponseData(output *usecase.CreateRelationshipOutput) (FollowResponseData, error) {
	users := make([]UserResponseData, 0, len(output.Users))
	for _, uwr := range output.Users {
		resp, err := NewUserResponseData(uwr.User, uwr.FollowerIDs, uwr.FollowingIDs, uwr.Courses, uwr.Reviews)
		if err != nil {
			return FollowResponseData{}, err
		}
		users = append(users, resp)
	}

	currentUser, err := NewUserWithRelationsResponse(output.CurrentUser)
	if err != nil {
		return FollowResponseData{}, err
	}

	followedUser, err := NewUserWithRelationsResponse(output.FollowedUser)
	if err != nil {
		return FollowResponseData{}, err
	}

	return FollowResponseData{
		Users:        users,
		CurrentUser:  currentUser,
		FollowedUser: followedUser,
	}, nil
}
