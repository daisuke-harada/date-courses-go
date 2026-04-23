# Copilot Instructions

## プロジェクト概要

Go 製のデートコース提案 REST API。Echo v4 + GORM + PostgreSQL 構成。
OpenAPI スキーマ（`api/`）から `oapi-codegen` でサーバーインターフェースと型を自動生成し、カスタムジェネレータ（`handler_generator.go`）でハンドラースタブも生成する。

## 開発プロセス（Step-by-Step）

- **分割統治**: 一度に巨大な機能を実装せず、動作可能な最小単位（1機能、1関数）で提案すること。
- **合意形成**: 実装を始める前に、まず「これから行う実装のステップ」を箇条書きで提示し、私の承認を得ること。
- **インクリメンタル**: 各ステップが完了するたびに、ビルドとテストが通る状態を維持すること。
- **TDD**: 実装コードを書く前にテストを先に書く（Red → Green → Refactor）。テストなしの実装 PR はマージ禁止。
- **コード品質**: コードは読みやすく、保守しやすいものにすること。冗長なコードや複雑なロジックは避けること。
- **ビルド保証**: 全タスクが完了したら `make gen && go build ./...` を1回実行し、ビルドエラーがゼロになるまで修正を続けること。エラーが出た場合は該当箇所を調査して修正し、再度 `go build ./...` を実行する。ビルドが通らない状態で作業を終了しない。

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
│   ├── repository/                  # repository インターフェース定義
│   │   └── mock/                   # repository mock（package repomock）自動生成
│   └── service/                     # ドメインサービス インターフェース定義
│       └── mock/                   # service mock（package servicemock）自動生成
├── usecase/                         # ビジネスロジック（Input/Output 型定義）
│   └── mock/                       # usecase InputPort mock（package usecasemock）自動生成
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

## ディレクトリごとの実装方針

各ディレクトリの詳細な実装方針は、それぞれの `CLAUDE.md` を参照してください。

| ディレクトリ | 役割 | CLAUDE.md |
|---|---|---|
| `internal/domain/model/` | GORM タグ付き struct・enum 型定義 | [CLAUDE.md](internal/domain/model/CLAUDE.md) |
| `internal/domain/repository/` | repository インターフェース定義 | [CLAUDE.md](internal/domain/repository/CLAUDE.md) |
| `internal/domain/service/` | ドメインサービス インターフェース定義 | [CLAUDE.md](internal/domain/service/CLAUDE.md) |
| `internal/usecase/` | ビジネスロジック・バリデーション | [CLAUDE.md](internal/usecase/CLAUDE.md) |
| `internal/infrastructure/persistence/` | repository の GORM 実装 | [CLAUDE.md](internal/infrastructure/persistence/CLAUDE.md) |
| `internal/interface/handler/` | Echo ハンドラー実装・型変換 | [CLAUDE.md](internal/interface/handler/CLAUDE.md) |
| `internal/interface/openapi/` | レスポンス変換関数（編集禁止ファイルあり） | [CLAUDE.md](internal/interface/openapi/CLAUDE.md) |

## ビルド・開発コマンド

makefile に記載。主なコマンド：

```sh
make deps           # go mod download
make docker-up      # Docker 起動 + PostgreSQL 待機
make apply-schema   # psqldef でスキーマ適用（差分のみ）
make db-seed        # シードデータ投入
make db-drop        # スキーマ全削除（NOT NULL 追加時などに使用）
make gen            # openapi-generate + go-generate（全タスク完了後に go build ./... と合わせて実行）
make run            # サーバー起動
```

> **注意**: `make gen` と `go build ./...` は全タスク完了後に **1回だけ** 実行する。
> エラーが出た場合は該当箇所を調査して修正し、再度 `go build ./...` を実行する。ビルドが通らない状態で作業を終了しない。

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

### apperror 関数の使い方

各エラー関数には「cause なし版」と「WithCause 版」の2種類があります。

| 用途 | 関数 |
|---|---|
| cause なし（クライアントエラーのみ） | `apperror.NotFound()` / `apperror.BadRequest()` など |
| cause あり（slog にログを残したい） | `apperror.NotFoundWithCause(err)` / `apperror.UnprocessableEntityWithCause(err)` など |

```go
// ❌ Bad: error 型を string 引数の関数に直接渡せない
return apperror.UnprocessableEntity(err)

// ✅ Good: WithCause 版で err を slog 用にラップ（クライアントにはデフォルトメッセージ）
return apperror.UnprocessableEntityWithCause(err)

// ✅ Good: メッセージもカスタマイズしたい場合
return apperror.NotFoundWithCause(err, "ユーザーが見つかりません")

// ✅ Good: 500 は cause 必須（もともと cause を渡す設計）
return apperror.InternalServerError(err)
```

slog への出力は `error_handler.go` の `CustomHTTPErrorHandler` が自動的に行います。  
cause があれば `slog.Error`、cause なしの 4xx は `slog.Warn` として記録されます。

## TDD（テスト駆動開発）

このプロジェクトでは TDD（テスト駆動開発）を採用しています。必ず以下のサイクルで開発を進めてください。

### Red → Green → Refactor サイクル

1. **Red**: 失敗するテストを先に書く（実装コードはまだ書かない）
2. **Green**: テストが通る最小限の実装を行う
3. **Refactor**: 動作を維持しながらコードを整理・改善する

### 複数 API を同時実装する場合（Branch Chaining 方式）

複数の API を同時に実装する際は、**branch chaining** パターンを採用してください。これにより、各機能の PR が独立しつつ、依存関係を管理できます。

#### Branch Chaining の手順

1. **最初の branch を作成**: `feature/issue-N1-summary` を `main` から作成
2. **最初の PR を作成**: 完了後、PR を作成（まだマージしない）
3. **次の branch を作成**: `feature/issue-N2-summary` を `feature/issue-N1-summary` から作成
4. **次の PR を作成**: 完了後、PR を作成（ベースブランチは `feature/issue-N1-summary`）
5. **以降繰り返し**: N3, N4, ... と続ける

#### マージ順序

上流から順にマージする：
1. issue-N1 をマージ
2. issue-N2 の base を `main` に変更してマージ
3. 以降同様

### 実装順序

新しい API を実装する際は、必ず以下の順序で進める：

| ステップ | 対象 | 内容 |
|---|---|---|
| 1 | usecase テスト | `XxxInteractor.Execute` の正常系・異常系を先に記述 |
| 2 | usecase 実装 | テストが通るように `Interactor` を実装 |
| 3 | usecase mock | `internal/usecase/mock/` に `MockXxxInputPort` を追加（`make gen` または手動） |
| 4 | handler テスト | `XxxHandler.Xxx` の正常系・バリデーション異常系・usecase エラー系を先に記述 |
| 5 | handler 実装 | テストが通るように `Handler` を実装 |
| 6 | DI 登録 | `di/infrastructure.go` に usecase を追加 |

### PR マージ条件

- usecase テスト・handler テストが両方存在すること
- `go test ./...` がすべて通ること
- テストなしの実装 PR はマージ禁止

## テスト規約

### パッケージ
- usecase テスト: `package usecase_test`（`internal/usecase/` 直下に配置）
- handler テスト: `package handler_test`（`internal/interface/handler/` 直下に配置）

### テスト形式
- **必ずサブテスト形式**（`t.Run`）を使う
- テスト関数名・サブテスト名は **英語スネークケース**（日本語・日本語混じり禁止）
- 命名規則: `TestXxx_<subtest_name>`

```go
// ✅ Good
func TestCreateRelationshipInteractor_Execute(t *testing.T) {
    t.Run("success", func(t *testing.T) { ... })
    t.Run("error_follow_self", func(t *testing.T) { ... })
    t.Run("error_current_user_not_found", func(t *testing.T) { ... })
}
```

### モックの使い方
モックは `go.uber.org/mock/gomock` を使用する。

| 対象 | 配置場所 | パッケージ名 | import alias |
|---|---|---|---|
| repository interface | `internal/domain/repository/mock/` | `repomock` | `repomock` |
| service interface | `internal/domain/service/mock/` | `servicemock` | `servicemock` |
| usecase InputPort | `internal/usecase/mock/` | `usecasemock` | `usecasemock` |

### handler テストのセットアップ
```go
func setupFormRequest(method, path string, form url.Values) (echo.Context, *httptest.ResponseRecorder) {
    e := echo.New()
    req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
    req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
    rec := httptest.NewRecorder()
    return e.NewContext(req, rec), rec
}
```

### アサーション
- `require.NoError` / `require.Error` → 以降の処理が前提条件に依存する場合
- `assert.Equal` / `assert.Contains` → 個別の値検証
- `ctrl.Finish()` を `defer` で呼ぶことで期待呼び出しを自動検証する

## 環境変数

`.envrc`（direnv）で管理。`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME` が必要。
