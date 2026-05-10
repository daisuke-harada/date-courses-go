# date-course-OpenAPI
バックエンド用の OpenAPI リポジトリ

## アーキテクチャ概要

### モノリス構成（Echo）

本プロジェクトは [Echo](https://echo.labstack.com/) を使用したモノリス構成を採用しています。
Echo は Lambda 上でも `github.com/awslabs/aws-lambda-go-api-proxy` 経由でそのまま動作するため、開発環境・Lambda 環境の両方で同一コードベースを使い回せます。マイクロサービス化に向けた段階的な分割が必要になった場合は、`cmd/lambda/` 配下に API 単位のエントリポイントを追加する設計にしています。

### AWS Lambda デプロイ（コスト削減 + 学習目的）

本番環境は AWS Lambda（arm64 / `provided.al2023` ランタイム）+ API Gateway HTTP API で稼働します。  
Lambda を採用した理由は次の 2 点です。

- **コスト削減**: 常時稼働のサーバが不要なため、リクエストがない時間帯の費用がゼロになります。
- **学習目的**: Go × SAM × Lambda の構成を実践的に習得するための題材として位置づけています。

データベースは TiDB Cloud（外部エンドポイント）を使用しており、VPC 設定は不要です。シークレット（DB 接続情報・JWT・Google Maps API キー）は AWS SSM Parameter Store で管理します。

---

## ディレクトリ構成（主要部分）

```
cmd/
  api/          # ローカル開発用サーバエントリポイント (go run)
  lambda/
    monolith/   # Lambda 用エントリポイント（全ルートをまとめたモノリス）
template.yaml   # AWS SAM テンプレート（Lambda + HTTP API の IaC 定義）
.github/
  workflows/
    ci.yaml     # CI: PR 時にビルド・Lint・テストを実行
    cd.yaml     # CD: main マージ後に SAM Deploy を実行
```

---

## CI/CD

### CI（`.github/workflows/ci.yaml`）

Pull Request をトリガーに次のジョブを実行します。

| ステップ | 内容 |
|---|---|
| Build | `go build ./...` |
| Lint | `golangci-lint` |
| Test | `go test ./...` |

同一 PR への連続プッシュは `concurrency` 設定により前のジョブをキャンセルします。

### CD（`.github/workflows/cd.yaml`）

`main` ブランチへのプッシュをトリガーに AWS へ自動デプロイします。

| ステップ | 内容 |
|---|---|
| OIDC 認証 | GitHub Actions OIDC で `AWS_DEPLOY_ROLE_ARN` のロールを一時取得 |
| SAM Build | `make sam-build`（arm64 向け Go クロスコンパイル + `sam build`） |
| SAM Deploy | `sam deploy` でスタック `date-courses-go` を `ap-northeast-1` へ反映 |

AWS 認証情報をシークレットに保存する必要はなく、OIDC により最小権限での一時認証を実現しています。

---

## SAM テンプレート（`template.yaml`）

`template.yaml` は AWS SAM（Serverless Application Model）の IaC 定義です。主なリソース構成は次のとおりです。

| リソース | 種別 | 説明 |
|---|---|---|
| `DateCoursesFunction` | `AWS::Serverless::Function` | arm64 / provided.al2023 のモノリス Lambda |
| `DateCoursesHttpApi` | `AWS::Serverless::HttpApi` | API Gateway HTTP API（CORS 設定済み） |
| `LambdaExecutionRole` | `AWS::IAM::Role` | Lambda 実行ロール（CloudWatch Logs のみ許可） |

環境変数はすべて SSM Parameter Store 参照（`AWS::SSM::Parameter::Value<String>`）で注入します。デプロイ時にシークレットがテンプレートにハードコードされることはありません。

### ローカルで Lambda を動かす

```bash
# 1. arm64 向けにビルドして SAM テンプレートをビルド
make sam-build

# 2. ローカル HTTP サーバとして起動（ポート 3000）
make local-lambda
# => sam local start-api
```

---

## apiの起動手順（ローカル開発）

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
go run ./cmd/api/main.go
```

必要な環境変数（例）:

```bash
export DB_USER=dev
export DB_PASSWORD=secret
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=date_courses_dev
export JWT_SECRET_KEY=your-secret
export GOOGLE_MAPS_API_KEY=your-key
```

（環境や開発ルールに合わせて適宜 `.env` や direnv、docker-compose の env を使ってください）

---

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
- oapi-codegen 実行（`openapi` パッケージ）
  - `internal/interface/openapi/generate.go` に `//go:generate` 指示がある
    - types: `api_types.gen.go`（`-generate types`）
    - echo-server: `api_server.gen.go`（`-generate echo-server`）
    - 生成後に handler ジェネレータを呼び出す: `//go:generate go run ../handler_generator.go`
- handler ジェネレータ（自作）
  - `internal/interface/handler_generator.go`
    - `openapi.ServerInterface` を reflect で読み、テンプレートから handler ファイルを生成
    - テンプレート: `templates/handler.tmpl`, `templates/handler_constructor.tmpl`
    - 出力先: `internal/interface/handler/`
    - 挙動: 既存ファイルが存在する場合は生成をスキップ（手書き実装を上書きしない）

### 実行手順（推奨）

```bash
make gen
```

### 運用ルール（重要）
- `internal/interface/openapi` は「生成物専用」フォルダにしてください。手書きコードを混ぜないこと。
- handler ジェネレータは既存ファイルを上書きしないため、消したくない実装はそのまま残ります。OpenAPI の変更で新しい operation が増えたら `go generate ./internal/interface/openapi` を実行して足りない handler を生成してください。
- generator 実行は `openapi` パッケージ経由で行う（`go generate ./internal/interface/openapi`）。`handler_generator.go` を直接 `go run` すると相対パスやビルド条件で期待通り動かないことがあります。

---

## データベーススキーマの適用（mysqldef）

このリポジトリではデータベーススキーマの適用に `mysqldef` (sqldef) を利用します。スキーマ定義ファイルは `internal/infrastructure/db/schema.sql` にあります。

### ローカル（Docker MySQL）

```bash
make apply-schema
```

### TiDB Cloud（本番相当）

```bash
make tidb-apply-schema
```

TiDB Cloud 接続時は `--ssl-mode=REQUIRED` が付与されます。

必要な環境変数:

```bash
export DB_USER=dev
export DB_PASSWORD=secret
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=date_courses_dev
```
