# date-course-OpenAPI
バックエンド用の OpenAPI リポジトリ

## apiの起動手順

前提:

- Docker / docker-compose が利用できること
- Go ツールチェイン（ビルドと `make run` に使用）

アプリの起動（開発環境の一例）:

1. 必要なコンテナを立ち上げます:

```bash
docker compose up -d
```

2. アプリを起動します（Makefile の `run` ターゲットを使用）:

```bash
make run
```

Makefile の `run` は次のコマンドを実行します:

```text
go run ./cmd/main.go
```

必要な環境変数（例）:

```bash
export POSTGRES_USER=dev
export POSTGRES_PASSWORD=secret
export POSTGRES_HOST=127.0.0.1
export POSTGRES_PORT=5432
export POSTGRES_DB=date_courses_dev
```

（環境や開発ルールに合わせて適宜 `.env` や direnv、docker-compose の env を使ってください）

## 自動生成について

このプロジェクトの自動生成フローをまとめます。現在の実装に基づく内容です。

### 概要
生成は次の順で行います。

1. OpenAPI を解決して `api/resolved/openapi/openapi.yaml` を作成（スクリプト使用）
2. `oapi-codegen` で型・Echo サーバインターフェースを生成（`gen` パッケージ）
3. 自前ジェネレータで `gen.ServerInterface` に合わせた handler スタブを生成（既存ファイルは上書きしない）

### 主要ファイル
- OpenAPI 解決スクリプト
  - `scripts/openapi-generator-cli.sh` — Docker で openapi-generator-cli を呼ぶ（出力先: `api/resolved`）
- oapi-codegen 実行（`gen` パッケージ）
  - `internal/infrastructure/cmd/api/gen/generate.go` に `//go:generate` 指示がある
    - types: `api_types.gen.go`（`-generate types`）
    - echo-server: `api_server.gen.go`（`-generate echo-server`）
    - 生成後に handler ジェネレータを呼び出す: `//go:generate go run ../handler_generator.go`
- handler ジェネレータ（自作）
  - `internal/infrastructure/cmd/api/handler_generator.go`
    - `gen.ServerInterface` を reflect で読み、テンプレートから handler ファイルを生成
    - テンプレート: `templates/handler.tmpl`, `templates/handler_constructor.tmpl`
    - 出力先: `internal/infrastructure/cmd/api/handler/`
    - 挙動: 既存ファイルが存在する場合は生成をスキップ（手書き実装を上書きしない）

### 実行手順（推奨）

```bash
make gen
```

### 運用ルール（重要）
- `internal/infrastructure/cmd/api/gen` は「生成物専用」フォルダにしてください。手書きコードを混ぜないこと。
- handler ジェネレータは既存ファイルを上書きしないため、消したくない実装はそのまま残ります。OpenAPI の変更で新しい operation が増えたら `go generate ./internal/infrastructure/cmd/api/gen` を実行して足りない handler を生成してください。
- generator 実行は `gen` パッケージ経由で行う（`go generate ./internal/infrastructure/cmd/api/gen`）。`handler_generator.go` を直接 `go run` すると相対パスやビルド条件で期待通り動かないことがあります。

※ 現在のジェネレータ実装（`internal/infrastructure/cmd/api/gen/generate.go` と `internal/infrastructure/cmd/api/handler_generator.go`、および `scripts/openapi-generator-cli.sh`）に基づいています。運用方針（既存ファイルをスキップする）は現在の設定どおりです。ファイル配置やジェネレータの挙動を変更する場合は README を更新してください。

## データベーススキーマの適用（sqldef / psqldef）

このリポジトリではデータベーススキーマの適用に `psqldef` (sqldef) を利用します。スキーマ定義ファイルは `internal/infrastructure/db/schema.sql` にあります。

ポイント:
- スキーマ適用は `psqldef` コマンドで行われます。Makefile の `apply-schema` ターゲットは以下のコマンドを実行します:

```bash
psqldef -U ${POSTGRES_USER} --password ${POSTGRES_PASSWORD}  -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} ${POSTGRES_DB} < ./internal/infrastructure/db/schema.sql
```

- 事前準備:
  - `psqldef` をインストールしてください。例:
    - Homebrew が使える環境:

```bash
brew install psqldef
# または (tap が必要な場合があります)
# brew tap rhysd/tap && brew install psqldef
```

    - Go ツールチェインからインストールする場合:

```bash
go install github.com/k0kubun/psqldef/cmd/psqldef@latest
```

  - 環境変数 `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_HOST` / `POSTGRES_PORT` / `POSTGRES_DB` を適切に設定してください。

- 適用方法:
  - 開発環境で手軽に適用するには `make apply-schema` を実行します（`make gen` は自動生成フローに加えて `apply-schema` を呼びます）:

```bash
# 例
export POSTGRES_USER=dev
export POSTGRES_PASSWORD=secret
export POSTGRES_HOST=127.0.0.1
export POSTGRES_PORT=5432
export POSTGRES_DB=date_courses_dev
make apply-schema
```

- 注意点:
  - `psqldef` はスキーマ差分を DB に適用するツールです。実行前に重要なデータのバックアップを取り、CI 環境では機密情報（パスワード等）をシークレットで管理してください。
  - `apply-schema` の具体的な挙動（テーブルの変更や削除など）については `psqldef` のドキュメントを参照してください。

- 参照ファイル:
  - スキーマ本体: `internal/infrastructure/db/schema.sql`

必要であれば README に CI 用の例（接続文字列を使う方法や、PGPASSWORD 環境変数を使う例など）を追加します。ご希望あれば教えてください。
