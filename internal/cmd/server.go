//go:generate oapi-codegen -generate types -o api/api_types.gen.go -package api ../../api/resolved/openapi/openapi.yaml
//go:generate oapi-codegen -generate echo-server -o api/api_server.gen.go -package api ../../api/resolved/openapi/openapi.yaml
//go:generate go run handler_generator.go

package api

import (
	"github.com/daisuke-harada/date-courses-go/internal/cmd/api"
	"github.com/daisuke-harada/date-courses-go/internal/cmd/handler"
	"github.com/labstack/echo"
	"go.uber.org/dig"
)

func Run() {
	container := dig.New()
	container.Provide(NewEcho)
	container.Provide(handler.NewHandler)
	container.Invoke(func(router api.EchoRouter, si api.ServerInterface) {
		api.RegisterHandlers(router, si)
	})
}

func NewEcho() *echo.Echo {
	e := echo.New()
	return e
}
