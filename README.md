# date-course-OpenAPI
バックエンド用のOpenAPIリポジトリ

# 自動生成について

このプロジェクトの自動生成フローをまとめます。現在の実装に基づく内容です。

## 概要
生成は次の順で行います。

1. OpenAPI を解決して `api/resolved/openapi/openapi.yaml` を作成（スクリプト使用）
2. `oapi-codegen` で型・Echo サーバインターフェースを生成（`gen` パッケージ）
3. 自前ジェネレータで `gen.ServerInterface` に合わせた handler スタブを生成（既存ファイルは上書きしない）

## 主要ファイル
- OpenAPI 解決スクリプト
  `scripts/openapi-generator-cli.sh` — Docker で openapi-generator-cli を呼ぶ（出力先: `api/resolved`）
- oapi-codegen 実行（`gen` パッケージ）
  `internal/infrastructure/cmd/api/gen/generate.go` に `//go:generate` 指示がある
  - types: `api_types.gen.go`（`-generate types`）
  - echo-server: `api_server.gen.go`（`-generate echo-server`）
  - 生成後に handler ジェネレータを呼び出す: `//go:generate go run ../handler_generator.go`
- handler ジェネレータ（自作）
  `internal/infrastructure/cmd/api/handler_generator.go`
  - `gen.ServerInterface` を reflect で読み、テンプレートから handler ファイルを生成
  - テンプレート: `templates/handler.tmpl`, `templates/handler_constructor.tmpl`
  - 出力先: `internal/infrastructure/cmd/api/handler/`
  - 挙動: 既存ファイルが存在する場合は生成をスキップ（手書き実装を上書きしない）

## 実行手順（推奨）

```
make gen
```

## 運用ルール（重要）
- `internal/infrastructure/cmd/api/gen` は「生成物専用」フォルダにしてください。手書きコードを混ぜないこと。
- handler ジェネレータは既存ファイルを上書きしないため、消したくない実装はそのまま残ります。OpenAPI の変更で新しい operation が増えたら `go generate ./internal/infrastructure/cmd/api/gen` を実行して足りない handler を生成してください。
- generator 実行は `gen` パッケージ経由で行う（`go generate ./internal/infrastructure/cmd/api/gen`）。`handler_generator.go` を直接 `go run` すると相対パスやビルド条件で期待通り動かないことがあります。

※ 現在のジェネレータ実装（`internal/infrastructure/cmd/api/gen/generate.go` と `internal/infrastructure/cmd/api/handler_generator.go`、および `scripts/openapi-generator-cli.sh`）に基づいています。運用方針（既存ファイルをスキップする）は現在の設定どおりです。ファイル配置やジェネレータの挙動を変更する場合は README を更新してください。
