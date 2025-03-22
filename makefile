gen: openapi-generate go-generate

openapi-generate:
	bash scripts/openapi-generator-cli.sh

go-generate:
	go generate ./internal/api
