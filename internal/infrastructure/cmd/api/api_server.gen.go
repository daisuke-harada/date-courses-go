// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /)
	Get(ctx echo.Context) error

	// (GET /api/v1/courses)
	GetApiV1Courses(ctx echo.Context) error

	// (POST /api/v1/courses)
	PostApiV1Courses(ctx echo.Context) error

	// (POST /api/v1/courses/sort)
	PostApiV1CoursesSort(ctx echo.Context) error

	// (DELETE /api/v1/courses/{id})
	DeleteApiV1CoursesId(ctx echo.Context, id int) error

	// (GET /api/v1/courses/{id})
	GetApiV1CoursesId(ctx echo.Context, id int) error

	// (POST /api/v1/date_spot_name_search)
	PostApiV1DateSpotNameSearch(ctx echo.Context) error

	// (POST /api/v1/date_spot_reviews)
	PostApiV1DateSpotReviews(ctx echo.Context) error

	// (DELETE /api/v1/date_spot_reviews/{id})
	DeleteApiV1DateSpotReviewsId(ctx echo.Context, id int) error

	// (PUT /api/v1/date_spot_reviews/{id})
	PutApiV1DateSpotReviewsId(ctx echo.Context, id int) error

	// (GET /api/v1/date_spots)
	GetApiV1DateSpots(ctx echo.Context) error

	// (POST /api/v1/date_spots)
	PostApiV1DateSpots(ctx echo.Context) error

	// (POST /api/v1/date_spots/sort)
	PostApiV1DateSpotsSort(ctx echo.Context) error

	// (DELETE /api/v1/date_spots/{id})
	DeleteApiV1DateSpotsId(ctx echo.Context, id int) error

	// (GET /api/v1/date_spots/{id})
	GetApiV1DateSpotsId(ctx echo.Context, id int) error

	// (PUT /api/v1/date_spots/{id})
	PutApiV1DateSpotsId(ctx echo.Context, id int) error

	// (GET /api/v1/genres/{id})
	GetApiV1GenresId(ctx echo.Context, id int) error

	// (POST /api/v1/login)
	PostApiV1Login(ctx echo.Context) error

	// (GET /api/v1/prefectures/{id})
	GetApiV1PrefecturesId(ctx echo.Context, id int) error

	// (POST /api/v1/relationships)
	PostApiV1Relationships(ctx echo.Context) error

	// (DELETE /api/v1/relationships/{current_user_id}/{other_user_id})
	DeleteApiV1RelationshipsCurrentUserIdOtherUserId(ctx echo.Context, currentUserId int, otherUserId int) error

	// (POST /api/v1/signup)
	PostApiV1Signup(ctx echo.Context) error

	// (GET /api/v1/top)
	GetApiV1Top(ctx echo.Context) error

	// (POST /api/v1/user_name_search)
	PostApiV1UserNameSearch(ctx echo.Context) error

	// (GET /api/v1/users)
	GetApiV1Users(ctx echo.Context) error

	// (DELETE /api/v1/users/{id})
	DeleteApiV1UsersId(ctx echo.Context, id int) error

	// (GET /api/v1/users/{id})
	GetApiV1UsersId(ctx echo.Context, id int) error

	// (PUT /api/v1/users/{id})
	PutApiV1UsersId(ctx echo.Context, id int) error

	// (GET /api/v1/users/{user_id}/followers)
	GetApiV1UsersUserIdFollowers(ctx echo.Context, userId int) error

	// (GET /api/v1/users/{user_id}/followings)
	GetApiV1UsersUserIdFollowings(ctx echo.Context, userId int) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Get converts echo context to params.
func (w *ServerInterfaceWrapper) Get(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Get(ctx)
	return err
}

// GetApiV1Courses converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1Courses(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1Courses(ctx)
	return err
}

// PostApiV1Courses converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1Courses(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1Courses(ctx)
	return err
}

// PostApiV1CoursesSort converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1CoursesSort(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1CoursesSort(ctx)
	return err
}

// DeleteApiV1CoursesId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteApiV1CoursesId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteApiV1CoursesId(ctx, id)
	return err
}

// GetApiV1CoursesId converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1CoursesId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1CoursesId(ctx, id)
	return err
}

// PostApiV1DateSpotNameSearch converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1DateSpotNameSearch(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1DateSpotNameSearch(ctx)
	return err
}

// PostApiV1DateSpotReviews converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1DateSpotReviews(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1DateSpotReviews(ctx)
	return err
}

// DeleteApiV1DateSpotReviewsId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteApiV1DateSpotReviewsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteApiV1DateSpotReviewsId(ctx, id)
	return err
}

// PutApiV1DateSpotReviewsId converts echo context to params.
func (w *ServerInterfaceWrapper) PutApiV1DateSpotReviewsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutApiV1DateSpotReviewsId(ctx, id)
	return err
}

// GetApiV1DateSpots converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1DateSpots(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1DateSpots(ctx)
	return err
}

// PostApiV1DateSpots converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1DateSpots(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1DateSpots(ctx)
	return err
}

// PostApiV1DateSpotsSort converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1DateSpotsSort(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1DateSpotsSort(ctx)
	return err
}

// DeleteApiV1DateSpotsId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteApiV1DateSpotsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteApiV1DateSpotsId(ctx, id)
	return err
}

// GetApiV1DateSpotsId converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1DateSpotsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1DateSpotsId(ctx, id)
	return err
}

// PutApiV1DateSpotsId converts echo context to params.
func (w *ServerInterfaceWrapper) PutApiV1DateSpotsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutApiV1DateSpotsId(ctx, id)
	return err
}

// GetApiV1GenresId converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1GenresId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1GenresId(ctx, id)
	return err
}

// PostApiV1Login converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1Login(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1Login(ctx)
	return err
}

// GetApiV1PrefecturesId converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1PrefecturesId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1PrefecturesId(ctx, id)
	return err
}

// PostApiV1Relationships converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1Relationships(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1Relationships(ctx)
	return err
}

// DeleteApiV1RelationshipsCurrentUserIdOtherUserId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteApiV1RelationshipsCurrentUserIdOtherUserId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "current_user_id" -------------
	var currentUserId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "current_user_id", runtime.ParamLocationPath, ctx.Param("current_user_id"), &currentUserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter current_user_id: %s", err))
	}

	// ------------- Path parameter "other_user_id" -------------
	var otherUserId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "other_user_id", runtime.ParamLocationPath, ctx.Param("other_user_id"), &otherUserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter other_user_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteApiV1RelationshipsCurrentUserIdOtherUserId(ctx, currentUserId, otherUserId)
	return err
}

// PostApiV1Signup converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1Signup(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1Signup(ctx)
	return err
}

// GetApiV1Top converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1Top(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1Top(ctx)
	return err
}

// PostApiV1UserNameSearch converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1UserNameSearch(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1UserNameSearch(ctx)
	return err
}

// GetApiV1Users converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1Users(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1Users(ctx)
	return err
}

// DeleteApiV1UsersId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteApiV1UsersId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteApiV1UsersId(ctx, id)
	return err
}

// GetApiV1UsersId converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1UsersId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1UsersId(ctx, id)
	return err
}

// PutApiV1UsersId converts echo context to params.
func (w *ServerInterfaceWrapper) PutApiV1UsersId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutApiV1UsersId(ctx, id)
	return err
}

// GetApiV1UsersUserIdFollowers converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1UsersUserIdFollowers(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "user_id" -------------
	var userId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "user_id", runtime.ParamLocationPath, ctx.Param("user_id"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter user_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1UsersUserIdFollowers(ctx, userId)
	return err
}

// GetApiV1UsersUserIdFollowings converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1UsersUserIdFollowings(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "user_id" -------------
	var userId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "user_id", runtime.ParamLocationPath, ctx.Param("user_id"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter user_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1UsersUserIdFollowings(ctx, userId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/", wrapper.Get)
	router.GET(baseURL+"/api/v1/courses", wrapper.GetApiV1Courses)
	router.POST(baseURL+"/api/v1/courses", wrapper.PostApiV1Courses)
	router.POST(baseURL+"/api/v1/courses/sort", wrapper.PostApiV1CoursesSort)
	router.DELETE(baseURL+"/api/v1/courses/:id", wrapper.DeleteApiV1CoursesId)
	router.GET(baseURL+"/api/v1/courses/:id", wrapper.GetApiV1CoursesId)
	router.POST(baseURL+"/api/v1/date_spot_name_search", wrapper.PostApiV1DateSpotNameSearch)
	router.POST(baseURL+"/api/v1/date_spot_reviews", wrapper.PostApiV1DateSpotReviews)
	router.DELETE(baseURL+"/api/v1/date_spot_reviews/:id", wrapper.DeleteApiV1DateSpotReviewsId)
	router.PUT(baseURL+"/api/v1/date_spot_reviews/:id", wrapper.PutApiV1DateSpotReviewsId)
	router.GET(baseURL+"/api/v1/date_spots", wrapper.GetApiV1DateSpots)
	router.POST(baseURL+"/api/v1/date_spots", wrapper.PostApiV1DateSpots)
	router.POST(baseURL+"/api/v1/date_spots/sort", wrapper.PostApiV1DateSpotsSort)
	router.DELETE(baseURL+"/api/v1/date_spots/:id", wrapper.DeleteApiV1DateSpotsId)
	router.GET(baseURL+"/api/v1/date_spots/:id", wrapper.GetApiV1DateSpotsId)
	router.PUT(baseURL+"/api/v1/date_spots/:id", wrapper.PutApiV1DateSpotsId)
	router.GET(baseURL+"/api/v1/genres/:id", wrapper.GetApiV1GenresId)
	router.POST(baseURL+"/api/v1/login", wrapper.PostApiV1Login)
	router.GET(baseURL+"/api/v1/prefectures/:id", wrapper.GetApiV1PrefecturesId)
	router.POST(baseURL+"/api/v1/relationships", wrapper.PostApiV1Relationships)
	router.DELETE(baseURL+"/api/v1/relationships/:current_user_id/:other_user_id", wrapper.DeleteApiV1RelationshipsCurrentUserIdOtherUserId)
	router.POST(baseURL+"/api/v1/signup", wrapper.PostApiV1Signup)
	router.GET(baseURL+"/api/v1/top", wrapper.GetApiV1Top)
	router.POST(baseURL+"/api/v1/user_name_search", wrapper.PostApiV1UserNameSearch)
	router.GET(baseURL+"/api/v1/users", wrapper.GetApiV1Users)
	router.DELETE(baseURL+"/api/v1/users/:id", wrapper.DeleteApiV1UsersId)
	router.GET(baseURL+"/api/v1/users/:id", wrapper.GetApiV1UsersId)
	router.PUT(baseURL+"/api/v1/users/:id", wrapper.PutApiV1UsersId)
	router.GET(baseURL+"/api/v1/users/:user_id/followers", wrapper.GetApiV1UsersUserIdFollowers)
	router.GET(baseURL+"/api/v1/users/:user_id/followings", wrapper.GetApiV1UsersUserIdFollowings)

}
