package openapi

//go:generate oapi-codegen -generate types -o api_types.gen.go -package openapi ../../../../../api/resolved/openapi/openapi.yaml
//go:generate oapi-codegen -generate echo-server -o api_server.gen.go -package openapi ../../../../../api/resolved/openapi/openapi.yaml
//go:generate go run ../handler_generator.go
