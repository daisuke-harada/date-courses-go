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

go-generate: mock-generate
	go generate ./internal/interface/openapi

mock-generate: mock-repository mock-service mock-usecase

mock-repository:
	mkdir -p internal/domain/repository/mock
	for f in internal/domain/repository/*.go; do \
		mockgen -source=$$f -destination=internal/domain/repository/mock/$$(basename $$f) -package=repositorymock; \
	done

mock-service:
	mkdir -p internal/domain/service/mock
	for f in internal/domain/service/*.go; do \
		mockgen -source=$$f -destination=internal/domain/service/mock/$$(basename $$f) -package=servicemock; \
	done

mock-usecase:
	mkdir -p internal/usecase/mock
	for f in internal/usecase/*.go; do \
		case $$f in *_test.go) continue;; esac; \
		mockgen -source=$$f -destination=internal/usecase/mock/$$(basename $$f) -package=usecasemock; \
	done

run:
	go run ./cmd/api/main.go

db-seed:
	go run ./tools/seed/main.go

db-drop:
	docker compose exec db psql -U "${DB_USER}" -d "${DB_NAME}" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
