//go:generate oapi-codegen -generate types -o api_types.gen.go -package api ../../api/resolved/openapi/openapi.yaml
//go:generate oapi-codegen -generate echo-server -o api_server.gen.go -package api ../../api/resolved/openapi/openapi.yaml
//go:generate go run ../../scripts/handler.gen.go

package api

// func Run() {
// 	container := dig.New()

// 	container.Provide(NewEcho)

// }

// func NewEcho() *echo.Echo {
// 	e := echo.New()
// 	return e
// }
