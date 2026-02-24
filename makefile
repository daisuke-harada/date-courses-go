gen: openapi-generate go-generate apply-schema

apply-schema:
	psqldef -U ${POSTGRES_USER} --password ${POSTGRES_PASSWORD}  -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} ${POSTGRES_DB} < ./internal/infrastructure/db/schema.sql

openapi-generate:
	bash scripts/openapi-generator-cli.sh

go-generate:
	# ensure modules are tidy before generating code
	go generate ./internal/infrastructure/cmd/api/gen

run:
	go run ./cmd/main.go
