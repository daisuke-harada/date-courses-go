package openapi

import "github.com/daisuke-harada/date-courses-go/internal/usecase"

// UnFollowResponseBody は DELETE /api/v1/relationships のレスポンス型です。
// UnFollowResponseData の互換型で、UserResponseBody を使います。
type UnFollowResponseBody struct {
	Users          []UserResponseBody `json:"users"`
	CurrentUser    UserResponseBody   `json:"current_user"`
	UnfollowedUser UserResponseBody   `json:"unfollowed_user"`
}

// BuildDeleteRelationshipResponse は DeleteRelationshipOutput から UnFollowResponseBody を構築します。
func BuildDeleteRelationshipResponse(output *usecase.DeleteRelationshipOutput) (UnFollowResponseBody, error) {
	users := make([]UserResponseBody, 0, len(output.Users))
	for _, uwr := range output.Users {
		resp, err := BuildUserResponseBody(uwr.User, uwr.FollowerIDs, uwr.FollowingIDs, uwr.Courses, uwr.Reviews)
		if err != nil {
			return UnFollowResponseBody{}, err
		}
		users = append(users, resp)
	}

	currentUser, err := NewUserWithRelationsResponse(output.CurrentUser)
	if err != nil {
		return UnFollowResponseBody{}, err
	}

	unfollowedUser, err := NewUserWithRelationsResponse(output.UnfollowedUser)
	if err != nil {
		return UnFollowResponseBody{}, err
	}

	return UnFollowResponseBody{
		Users:          users,
		CurrentUser:    currentUser,
		UnfollowedUser: unfollowedUser,
	}, nil
}

// CreateRelationshipResponseBody は POST /api/v1/relationships のレスポンス型です。
type CreateRelationshipResponseBody struct {
	Users        []UserResponseBody `json:"users"`
	CurrentUser  UserResponseBody   `json:"current_user"`
	FollowedUser UserResponseBody   `json:"followed_user"`
}

// BuildCreateRelationshipResponse は CreateRelationshipOutput から
// CreateRelationshipResponseBody を構築します。
func BuildCreateRelationshipResponse(output *usecase.CreateRelationshipOutput) (CreateRelationshipResponseBody, error) {
	users := make([]UserResponseBody, 0, len(output.Users))
	for _, uwr := range output.Users {
		resp, err := BuildUserResponseBody(uwr.User, uwr.FollowerIDs, uwr.FollowingIDs, uwr.Courses, uwr.Reviews)
		if err != nil {
			return CreateRelationshipResponseBody{}, err
		}
		users = append(users, resp)
	}

	currentUser, err := NewUserWithRelationsResponse(output.CurrentUser)
	if err != nil {
		return CreateRelationshipResponseBody{}, err
	}

	followedUser, err := NewUserWithRelationsResponse(output.FollowedUser)
	if err != nil {
		return CreateRelationshipResponseBody{}, err
	}

	return CreateRelationshipResponseBody{
		Users:        users,
		CurrentUser:  currentUser,
		FollowedUser: followedUser,
	}, nil
}
