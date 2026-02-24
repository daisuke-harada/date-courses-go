gen: openapi-generate go-generate apply-schema

apply-schema:
	# use DB_* variables (kept in .envrc) and pass password via PGPASSWORD to avoid
	# treating an empty -U value as the next flag. This also avoids interactive prompt.
	PGPASSWORD="${DB_PASSWORD}" psqldef -U "${DB_USER}" -h "${DB_HOST}" -p "${DB_PORT}" "${DB_NAME}" < ./internal/infrastructure/db/schema.sql

openapi-generate:
	bash scripts/openapi-generator-cli.sh

go-generate:
	# ensure modules are tidy before generating code
	go generate ./internal/infrastructure/cmd/api/gen

run:
	go run ./cmd/main.go
