# set -o pipefail を使うため bash を明示する（デフォルトの sh では動かない環境がある）
SHELL := /bin/bash

setup: deps gen docker-up apply-schema db-seed

deps:
	go mod download

docker-up:
	docker compose up -d
	# wait for MySQL to be ready before proceeding
	docker compose exec db bash -c 'until mysqladmin ping -u "${DB_USER}" -p"${DB_PASSWORD}" --silent 2>/dev/null; do sleep 1; done'

gen: openapi-generate go-generate

apply-schema:
	mysqldef -u "${DB_USER}" -p "${DB_PASSWORD}" -h "${DB_HOST}" -P "${DB_PORT}" "${DB_NAME}" < ./internal/infrastructure/db/schema.sql

tidb-apply-schema:
	mysqldef -u "${DB_USER}" -p "${DB_PASSWORD}" -h "${DB_HOST}" -P "${DB_PORT}" --ssl-mode=REQUIRED "${DB_NAME}" < ./internal/infrastructure/db/schema.sql

tidb-seed:
	go run ./tools/seed/main.go

tidb-setup: tidb-apply-schema tidb-seed

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

# # Lambda ビルド: provided.al2023 ランタイム向け (arm64)
# # モノリス Lambda（全ルート）。API 個別化時は build-lambda-<api-name> を追加する。
# build-lambda-monolith:
# 	mkdir -p dist/lambda/monolith
# 	GOARCH=arm64 GOOS=linux go build -o dist/lambda/monolith/bootstrap ./cmd/lambda/monolith
#
# # SAM CLI でローカル実行 (template.yaml が必要)
# local-lambda:
# 	sam local start-api
#
# sam-build: build-lambda-monolith
# 	sam build
#
# sam-deploy: sam-build
# 	sam deploy \
# 		--stack-name date-courses-go \
# 		--region ap-northeast-1 \
# 		--no-confirm-changeset \
# 		--no-fail-on-empty-changeset \
# 		--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
# 		--resolve-s3
#
# sam-delete:
# 	aws cloudformation delete-stack \
# 		--stack-name date-courses-go \
# 		--region ap-northeast-1

# テスト出力から読む必要のない行を落とすフィルタ。
# テストを持たないパッケージ、対象外パッケージの警告、実行開始ログ（=== RUN）を除く。
TEST_FILTER := grep -v -e "no test files" -e "no tests to run" -e "^testing: warning" -e "^=== " -e "^PASS$$"

# テスト実行。パイプで grep を挟むため、go test の失敗を取りこぼさないよう pipefail を付ける。
test:
	@set -o pipefail; go test ./... | $(TEST_FILTER)

# サブテスト名まで表示する。何がテストされているか確認したいときに使う。
#   例: make test-v
test-v:
	@set -o pipefail; go test -v ./... | $(TEST_FILTER)

# 特定のテストだけ実行する。
#   例: make test-run RUN=TestGetCourseInteractor_Execute
#       make test-run RUN='TestCourseRepository.*'
test-run:
	@set -o pipefail; go test -v -run '$(RUN)' ./... | $(TEST_FILTER)

# キャッシュを使わず必ず実行する。
test-fresh:
	@set -o pipefail; go test -count=1 ./... | $(TEST_FILTER)

lint:
	golangci-lint run ./...

run:
	go run ./cmd/api/main.go

db-seed:
	go run ./tools/seed/main.go

db-drop:
	docker compose exec db mysql -u "${DB_USER}" -p"${DB_PASSWORD}" -e "DROP DATABASE IF EXISTS ${DB_NAME}; CREATE DATABASE ${DB_NAME};"


