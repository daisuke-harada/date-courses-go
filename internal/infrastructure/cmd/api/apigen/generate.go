package apigen

//go:generate oapi-codegen -generate types -o api_types.gen.go -package apigen ../../../../../api/resolved/openapi/openapi.yaml
//go:generate oapi-codegen -generate echo-server -o api_server.gen.go -package apigen ../../../../../api/resolved/openapi/openapi.yaml
//go:generate go run ../handler_generator.go
