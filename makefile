setup: deps gen docker-up apply-schema db-seed

deps:
	go mod download

docker-up:
	docker compose up -d
	# wait for PostgreSQL to be ready before proceeding
	docker compose exec db sh -c 'until pg_isready -U "$$POSTGRES_USER"; do sleep 1; done'

gen: openapi-generate go-generate

apply-schema:
	# use DB_* variables (kept in .envrc) and pass password via PGPASSWORD to avoid
	# treating an empty -U value as the next flag. This also avoids interactive prompt.
	PGPASSWORD="${DB_PASSWORD}" psqldef -U "${DB_USER}" -h "${DB_HOST}" -p "${DB_PORT}" "${DB_NAME}" < ./internal/infrastructure/db/schema.sql

openapi-generate:
	bash scripts/openapi-generator-cli.sh

go-generate:
	# ensure modules are tidy before generating code
	go generate ./internal/interface/openapi

run:
	go run ./cmd/api/main.go

db-seed:
	go run ./tools/seed/main.go

db-drop:
	docker compose exec db psql -U "${DB_USER}" -d "${DB_NAME}" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
