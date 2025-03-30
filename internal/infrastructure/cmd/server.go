//go:generate oapi-codegen -generate types -o api/api_types.gen.go -package api ../../../api/resolved/openapi/openapi.yaml
//go:generate oapi-codegen -generate echo-server -o api/api_server.gen.go -package api ../../../api/resolved/openapi/openapi.yaml
//go:generate go run handler_generator.go

package api

import (
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/handler"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/middleware"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func Run(log logger.LoggerInterface) error {
	container := dig.New()
	container.Provide(NewEcho)
	container.Provide(log)
	container.Provide(middleware.NewMiddlewre)
	container.Provide(handler.NewHandler)
	container.Invoke(func(e *echo.Echo, router api.EchoRouter, si api.ServerInterface, mw middleware.Middleware) error {
		e.Use(mw.LoggerMiddleware())
		api.RegisterHandlers(router, si)
		return e.Start(":8080")
	})
	return nil
}

func NewEcho() *echo.Echo {
	e := echo.New()
	return e
}
