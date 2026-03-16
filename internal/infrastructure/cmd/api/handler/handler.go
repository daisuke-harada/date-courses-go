package handler

import "github.com/daisuke-harada/date-courses-go/internal/usecase"

type Handler struct {
	DeleteApiV1CoursesIdHandler
	DeleteApiV1DateSpotReviewsIdHandler
	DeleteApiV1DateSpotsIdHandler
	DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler
	DeleteApiV1UsersIdHandler
	GetHandler
	GetApiV1CoursesHandler
	GetApiV1CoursesIdHandler
	GetApiV1DateSpotsHandler
	GetApiV1DateSpotsIdHandler
	GetApiV1GenresIdHandler
	GetApiV1PrefecturesIdHandler
	GetApiV1TopHandler
	GetApiV1UsersHandler
	GetApiV1UsersIdHandler
	GetApiV1UsersUserIdFollowersHandler
	GetApiV1UsersUserIdFollowingsHandler
	PostApiV1CoursesHandler
	PostApiV1CoursesSortHandler
	PostApiV1DateSpotNameSearchHandler
	PostApiV1DateSpotReviewsHandler
	PostApiV1DateSpotsHandler
	PostApiV1DateSpotsSortHandler
	PostApiV1LoginHandler
	PostApiV1RelationshipsHandler
	PostApiV1SignupHandler
	PostApiV1UserNameSearchHandler
	PutApiV1DateSpotReviewsIdHandler
	PutApiV1DateSpotsIdHandler
	PutApiV1UsersIdHandler
}

// dig が各 InputPort を解決して注入します。
// 新しいユースケースを追加する際は引数を追加するだけです。
func NewHandler(
	getDateSpotsInputPort usecase.GetDateSpotsInputPort,
) *Handler {
	return &Handler{
		DeleteApiV1CoursesIdHandler: DeleteApiV1CoursesIdHandler{},
		GetApiV1DateSpotsHandler: GetApiV1DateSpotsHandler{
			InputPort: getDateSpotsInputPort,
		},
	}
}
