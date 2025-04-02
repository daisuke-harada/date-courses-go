//go:generate oapi-codegen -generate types -o api_types.gen.go -package api ../../../../api/resolved/openapi/openapi.yaml
//go:generate oapi-codegen -generate echo-server -o api_server.gen.go -package api ../../../../api/resolved/openapi/openapi.yaml
//go:generate go run handler_generator.go

package api

import (
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api/handler"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/cmd/api/middleware"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func Run(log logger.LoggerInterface) error {
	container := dig.New()
	container.Provide(func() logger.LoggerInterface {
		return log
	})
	container.Provide(NewEcho)
	container.Provide(handler.NewHandler)
	container.Provide(middleware.NewMiddleware)
	err := container.Invoke(func(e *echo.Echo, ha *handler.Handler, mw middleware.MiddlewareInterface) error {
		e.Use(mw.Logger())
		e.Use(mw.Cors())
		RegisterHandlers(e, ha)
		if err := e.Start(":8080"); err != nil {
			return err
		}
		return nil
	})

	return err
}

func NewEcho() *echo.Echo {
	e := echo.New()
	return e
}
