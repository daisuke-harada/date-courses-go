# Copilot Instructions

## プロジェクト概要

Go 製のデートコース提案 REST API。Echo v4 + GORM + PostgreSQL 構成。
OpenAPI スキーマ（`api/`）から `oapi-codegen` でサーバーインターフェースと型を自動生成し、カスタムジェネレータ（`handler_generator.go`）でハンドラースタブも生成する。

## 開発プロセス（Step-by-Step）

- **分割統治**: 一度に巨大な機能を実装せず、動作可能な最小単位（1機能、1関数）で提案すること。
- **合意形成**: 実装を始める前に、まず「これから行う実装のステップ」を箇条書きで提示し、私の承認を得ること。
- **インクリメンタル**: 各ステップが完了するたびに、ビルドとテストが通る状態を維持すること。
- **コード品質**: コードは読みやすく、保守しやすいものにすること。冗長なコードや複雑なロジックは避けること。

## アーキテクチャ

クリーンアーキテクチャに基づいた以下の層構造：

```
handler (infrastructure) → usecase → repository (interface) → persistence (infrastructure)
```

## ディレクトリ構成

```
cmd/
└── api/
    └── main.go                      # エントリーポイント
api/
├── OpenAPI.yaml                     # OpenAPI スキーマ（$ref で分割管理）
├── paths/                           # パス定義
├── components/schemas/              # スキーマ定義
└── resolved/openapi/openapi.yaml    # 解決済みスキーマ（生成物）
internal/
├── apperror/
│   └── errors.go                    # アプリエラー型（NotFound, BadRequest など）
├── config/
│   └── config.go                    # 環境変数読み込み
├── di/
│   ├── container.go                 # DIコンテナ初期化
│   ├── di.go                        # usecase の provide
│   └── infrastructure.go            # persistence の provide
├── domain/
│   ├── model/                       # ドメインモデル（GORM タグ付き struct）
│   └── repository/                  # repository インターフェース定義
├── usecase/                         # ビジネスロジック（Input/Output 型定義）
├── infrastructure/
│   ├── db/
│   │   ├── db.go                    # GORM 接続
│   │   └── schema.sql               # テーブル定義（psqldef で適用）
│   └── persistence/                 # repository インターフェースの GORM 実装
└── interface/                        # HTTP 層 / 生成ツール類
    ├── server.go                    # Echo サーバーセットアップ (package iface)
    ├── handler_generator.go         # ハンドラースタブ自動生成ツール
    ├── handler/                     # Echo ハンドラー実装
    ├── middleware/                  # Echo ミドルウェア（error_handler, cors など）
    └── openapi/                     # oapi-codegen 生成ファイル（編集禁止）
pkg/
└── logger/                          # slog ラッパー
tools/
└── seed/main.go                     # シードデータ投入
```

## コーディングルール

### domain/model

- GORM タグ付きの struct（エンティティ）はここに定義する
- 複数のモデルをまとめた集約型（例: `XxxWithRelations`）もここに定義する
- usecase や openapi など他のパッケージの型には依存しない

### handler

- `openapi` パッケージの型（`openapi.XxxParams`, `openapi.XxxJSONRequestBody` など）を受け取り、usecase の `Input` 型に変換して渡す
- エラーは `return err` のみ。`ctx.JSON` でエラーを直接返さない（`CustomHTTPErrorHandler` が処理する）
- usecase が返した `error` はそのまま `return err` する
- XxxHandler の Xxx メソッドはHTTPメソッドになります

```go
func (h *XxxHandler) Xxx(ctx echo.Context, params openapi.XxxParams) error {
    input := usecase.XxxInput{Field: params.Field}
    output, err := h.InputPort.Execute(ctx.Request().Context(), input)
    if err != nil {
        return err
    }
    response, err := openapi.NewXxxResponse(output)
    if err != nil {
        return apperror.InternalServerError(err)
    }
    return ctx.JSON(http.StatusOK, response)
}
```

### usecase

- `openapi` パッケージに依存しない。専用の `Input`/`Output` 型を定義する
- `Output` 型はポインタ（`*XxxOutput`）で返す
- エラーは `apperror` パッケージを使って返す（`apperror.NotFound()`, `apperror.InternalServerError(err)` など）
- repository の `SearchParams` 型に変換して渡す
- XxxHandler の Xxx メソッドはHTTPメソッドになります

```go
type XxxInput struct { ... }
type XxxOutput struct { ... }

func (i *XxxInteractor) Execute(ctx context.Context, input XxxInput) (*XxxOutput, error) {
    result, err := i.Repo.Search(ctx, repository.XxxSearchParams{...})
    if err != nil {
        return nil, apperror.InternalServerError(err)
    }
    return &XxxOutput{...}, nil
}
```

### repository

- インターフェースには実際に使用するメソッドのみ定義する（未使用の `GetByID`/`Update`/`Delete` は定義しない）
- 検索条件は `XxxSearchParams` 型で受け取る（GORM に依存しない形）
- persistence 層でのみ GORM の WHERE 句を組み立てる
- ログは `slog.ErrorContext` / `slog.InfoContext` を使う。`fmt.Print` 系は禁止

```go
// repository インターフェース
type XxxRepository interface {
    Create(ctx context.Context, xxx *model.Xxx) error
    Search(ctx context.Context, params XxxSearchParams) ([]*model.Xxx, error)
}

// persistence 実装
func (r *xxxRepository) Search(ctx context.Context, params repository.XxxSearchParams) ([]*model.Xxx, error) {
    db := r.db.WithContext(ctx).Model(&model.Xxx{})
    if params.Name != nil {
        db = db.Where("name LIKE ?", "%"+*params.Name+"%")
    }
    ...
}
```

## ビルド・開発コマンド

makefile に記載。主なコマンド：

```sh
make deps           # go mod download
make docker-up      # Docker 起動 + PostgreSQL 待機
make apply-schema   # psqldef でスキーマ適用（差分のみ）
make db-seed        # シードデータ投入
make db-drop        # スキーマ全削除（NOT NULL 追加時などに使用）
make gen            # openapi-generate + go-generate
make run            # サーバー起動
```

> **注意**: `apply-schema` は既存データがある状態で `NOT NULL` カラムを追加するとエラーになる。
> その場合は `make db-drop && make apply-schema && make db-seed` の順で実行する。

## OpenAPI / コード生成

- スキーマ: `api/OpenAPI.yaml`（`$ref` で `api/paths/` と `api/components/` を参照）
- 生成ファイル: `internal/interface/openapi/api_server.gen.go`, `api_types.gen.go`（編集禁止）
- ハンドラースタブ生成: `go generate ./internal/interface/openapi`（既存ファイルは上書きしない）
- 全パスの `responses:` に `default:` エラーレスポンスを追加済み

## エラーハンドリング

- `internal/apperror/errors.go` にエラー型を定義
- handler から `return apperror.NotFound()` のように `error` 型で返す
- `internal/interface/middleware/error_handler.go` の `CustomHTTPErrorHandler` が受け取り JSON レスポンスを返す
- レスポンス形式: `{ "errorMessages": ["メッセージ"] }`

## 環境変数

`.envrc`（direnv）で管理。`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME` が必要。
