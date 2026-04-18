package handler

import (
	"github.com/daisuke-harada/date-courses-go/internal/di"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
)

func NewHandler(container *di.Container) *Handler {
	return &Handler{
		DeleteApiV1CoursesIdHandler: DeleteApiV1CoursesIdHandler{},
		DeleteApiV1DateSpotReviewsIdHandler: DeleteApiV1DateSpotReviewsIdHandler{
			InputPort: di.MustInvoke[usecase.DeleteDateSpotReviewInputPort](container),
		},
		DeleteApiV1DateSpotsIdHandler: DeleteApiV1DateSpotsIdHandler{
			InputPort: di.MustInvoke[usecase.DeleteDateSpotInputPort](container),
		},
		DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler: DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler{
			InputPort: di.MustInvoke[usecase.DeleteRelationshipInputPort](container),
		},
		DeleteApiV1UsersIdHandler: DeleteApiV1UsersIdHandler{
			InputPort: di.MustInvoke[usecase.DeleteUserInputPort](container),
		},
		GetHandler: GetHandler{},
		GetApiV1CoursesHandler: GetApiV1CoursesHandler{
			InputPort: di.MustInvoke[usecase.GetCoursesInputPort](container),
		},
		GetApiV1CoursesIdHandler: GetApiV1CoursesIdHandler{
			InputPort: di.MustInvoke[usecase.GetCourseInputPort](container),
		},
		GetApiV1DateSpotsHandler: GetApiV1DateSpotsHandler{
			InputPort: di.MustInvoke[usecase.GetDateSpotsInputPort](container),
		},
		GetApiV1DateSpotsIdHandler:   GetApiV1DateSpotsIdHandler{},
		GetApiV1GenresIdHandler:      GetApiV1GenresIdHandler{},
		GetApiV1PrefecturesIdHandler: GetApiV1PrefecturesIdHandler{},
		GetApiV1TopHandler:           GetApiV1TopHandler{},
		GetApiV1UsersHandler: GetApiV1UsersHandler{
			InputPort: di.MustInvoke[usecase.GetUsersInputPort](container),
		},
		GetApiV1UsersIdHandler: GetApiV1UsersIdHandler{
			InputPort: di.MustInvoke[usecase.GetUserInputPort](container),
		},
		GetApiV1UsersUserIdFollowersHandler: GetApiV1UsersUserIdFollowersHandler{
			InputPort: di.MustInvoke[usecase.GetUserFollowersInputPort](container),
		},
		GetApiV1UsersUserIdFollowingsHandler: GetApiV1UsersUserIdFollowingsHandler{
			InputPort: di.MustInvoke[usecase.GetUserFollowingsInputPort](container),
		},
		PostApiV1CoursesHandler: PostApiV1CoursesHandler{},
		PostApiV1DateSpotReviewsHandler: PostApiV1DateSpotReviewsHandler{
			InputPort: di.MustInvoke[usecase.CreateDateSpotReviewInputPort](container),
		},
		PostApiV1DateSpotsHandler: PostApiV1DateSpotsHandler{
			InputPort: di.MustInvoke[usecase.CreateDateSpotInputPort](container),
		},
		PostApiV1LoginHandler: PostApiV1LoginHandler{
			InputPort: di.MustInvoke[usecase.LoginInputPort](container),
		},
		PostApiV1RelationshipsHandler: PostApiV1RelationshipsHandler{
			InputPort: di.MustInvoke[usecase.CreateRelationshipInputPort](container),
		},
		PostApiV1SignupHandler: PostApiV1SignupHandler{
			InputPort: di.MustInvoke[usecase.SignupInputPort](container),
		},
		PutApiV1DateSpotReviewsIdHandler: PutApiV1DateSpotReviewsIdHandler{
			InputPort: di.MustInvoke[usecase.UpdateDateSpotReviewInputPort](container),
		},
		PutApiV1DateSpotsIdHandler: PutApiV1DateSpotsIdHandler{
			InputPort: di.MustInvoke[usecase.UpdateDateSpotInputPort](container),
		},
		PutApiV1UsersIdHandler: PutApiV1UsersIdHandler{
			InputPort: di.MustInvoke[usecase.UpdateUserInputPort](container),
		},
	}
}
