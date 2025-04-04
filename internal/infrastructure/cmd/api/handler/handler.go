package handler

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

func NewHandler() *Handler {
	return &Handler{}
}