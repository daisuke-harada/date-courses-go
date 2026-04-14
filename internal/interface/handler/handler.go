package handler

import (
	"github.com/daisuke-harada/date-courses-go/internal/di"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
)

func NewHandler(container *di.Container) *Handler {
	return &Handler{
		DeleteApiV1CoursesIdHandler:                             DeleteApiV1CoursesIdHandler{},
		DeleteApiV1DateSpotReviewsIdHandler:                     DeleteApiV1DateSpotReviewsIdHandler{},
		DeleteApiV1DateSpotsIdHandler:                           DeleteApiV1DateSpotsIdHandler{},
		DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler: DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler{},
		DeleteApiV1UsersIdHandler:                               DeleteApiV1UsersIdHandler{},
		GetHandler:                                              GetHandler{},
		GetApiV1CoursesHandler:                                  GetApiV1CoursesHandler{},
		GetApiV1CoursesIdHandler:                                GetApiV1CoursesIdHandler{},
		GetApiV1DateSpotsHandler: GetApiV1DateSpotsHandler{
			InputPort: di.MustInvoke[usecase.GetDateSpotsInputPort](container),
		},
		GetApiV1DateSpotsIdHandler: GetApiV1DateSpotsIdHandler{},
		GetApiV1GenresIdHandler:    GetApiV1GenresIdHandler{},
		GetApiV1PrefecturesIdHandler: GetApiV1PrefecturesIdHandler{},
		GetApiV1TopHandler:           GetApiV1TopHandler{},
		GetApiV1UsersHandler: GetApiV1UsersHandler{
			InputPort: di.MustInvoke[usecase.GetUsersInputPort](container),
		},
		GetApiV1UsersIdHandler: GetApiV1UsersIdHandler{
			InputPort: di.MustInvoke[usecase.GetUserInputPort](container),
		},
		GetApiV1UsersUserIdFollowersHandler:  GetApiV1UsersUserIdFollowersHandler{},
		GetApiV1UsersUserIdFollowingsHandler: GetApiV1UsersUserIdFollowingsHandler{},
		PostApiV1CoursesHandler:              PostApiV1CoursesHandler{},
		PostApiV1DateSpotReviewsHandler:      PostApiV1DateSpotReviewsHandler{},
		PostApiV1DateSpotsHandler:            PostApiV1DateSpotsHandler{},
		PostApiV1LoginHandler: PostApiV1LoginHandler{
			InputPort: di.MustInvoke[usecase.LoginInputPort](container),
		},
		PostApiV1RelationshipsHandler: PostApiV1RelationshipsHandler{},
		PostApiV1SignupHandler: PostApiV1SignupHandler{
			InputPort: di.MustInvoke[usecase.SignupInputPort](container),
		},
		PutApiV1DateSpotReviewsIdHandler: PutApiV1DateSpotReviewsIdHandler{},
		PutApiV1DateSpotsIdHandler:       PutApiV1DateSpotsIdHandler{},
		PutApiV1UsersIdHandler:           PutApiV1UsersIdHandler{},
	}
}
