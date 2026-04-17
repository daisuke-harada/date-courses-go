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

### `internal/domain/model/`
- GORM タグ付きの struct を定義する
- ビジネスロジックは持たない（値オブジェクトは型エイリアスで表現）
- 例: `type Gender string` + `const GenderMale Gender = "male"`

### ドメインモデルと列挙型（enum）に関する注意
- ドメインモデルに列挙型（enum）を追加する可能性があります。追加する場合は以下を守ってください。
    - 型は基本的に型エイリアスで定義する（例: `type Gender string`）。
    - 値は `const` で列挙する（例: `const ( GenderMale Gender = "male" ... )`）。
    - JSON エンコード/デコードや DB マイグレーションで特別な処理が必要な場合は、`MarshalJSON`/`UnmarshalJSON` を実装するか、変換ヘルパーを用意する。
    - OpenAPI 側に別の enum 型がある場合は、`internal/interface/openapi` 側に変換関数（`NewXxx` / `ToModelXxx`）を実装して一元管理すること。
    - 型の追加・変更を行ったら、該当するテストと `make gen`（openapi/types 再生成や mock 再生成）が必要になる点に注意する。


### `internal/domain/repository/`
- interface のみ定義。実装は `infrastructure/persistence/` に置く
- メソッドは実際に使うものだけ定義する（YAGNI）
- 検索条件は `XxxSearchParams` 型にまとめる
- `mock/` サブディレクトリに `package repomock` としてモックを自動生成する（`make gen`）

### `internal/domain/service/`
- 複数の repository をまたぐドメインロジックをインターフェースで定義
- `mock/` サブディレクトリに `package servicemock` としてモックを自動生成する（`make gen`）

### `internal/usecase/`
- `XxxInputPort` interface・`XxxInput`・`XxxOutput`・`XxxInteractor` を 1 ファイルにまとめる
- `openapi` パッケージへの依存禁止。handler 層で変換する
- `Output` はポインタ返し（`*XxxOutput`）
- エラーは必ず `apperror` パッケージ経由で返す
- `mock/` サブディレクトリに `package usecasemock` としてモックを自動生成する（`make gen`）

### `internal/infrastructure/persistence/`
- repository インターフェースの GORM 実装のみ
- GORM の WHERE 句はここだけに閉じ込める
- エラーは `apperror.InternalServerError(err)` か `apperror.NotFound()` に変換して返す

### `internal/interface/handler/`
- `XxxHandler` struct に `InputPort` フィールドを持つ
- フォームや JSON のパースはここで行い、usecase の `Input` 型に変換して渡す
- エラーは `return err` のみ（`ctx.JSON` でエラーを直接返さない）
- テストは同パッケージに `handler_test` パッケージで記述する

### `internal/interface/openapi/`
- `api_server.gen.go`, `api_types.gen.go` は **編集禁止**（`make gen` で再生成）
- レスポンス変換用の手書きファイル（`signup.go`, `login.go` など）はここに置く
- `Gender` 型など openapi 型への変換は `NewXxx()` 関数を定義する

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

#### Branch Chaining の例

```
main
 ├─ feature/issue-51-delete-date-spot (PR #64) ← 最初のブランチ
 │   ├─ feature/issue-52-update-date-spot (PR #65, base: issue-51) ← 2番目
 │   │   ├─ feature/issue-53-create-review (base: issue-52)
 │   │   └─ ...
```

#### マージ順序

上流から順にマージする：
1. issue-51 をマージ
2. issue-52 の base を `main` に変更してマージ
3. issue-53 の base を `main` に変更してマージ
4. 以降同様

#### メリット

- 複数 API の実装を並列化できる
- PR 単位で独立した変更を管理できる
- マージ時の競合を最小化できる
- 各 PR のレビューが小さく、容易

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

// ❌ Bad
func TestCreateRelationshipInteractor_Execute_正常系_フォロー成功(t *testing.T) { ... }
```

### モックの使い方
モックは `go.uber.org/mock/gomock` を使用する。

| 対象 | 配置場所 | パッケージ名 | import alias |
|---|---|---|---|
| repository interface | `internal/domain/repository/mock/` | `repomock` | `repomock` |
| service interface | `internal/domain/service/mock/` | `servicemock` | `servicemock` |
| usecase InputPort | `internal/usecase/mock/` | `usecasemock` | `usecasemock` |

```go
// usecase テストでの repository mock 使用例
import repomock "github.com/daisuke-harada/date-courses-go/internal/domain/repository/mock"

ctrl := gomock.NewController(t)
defer ctrl.Finish()
dateSpotRepo := repomock.NewMockDateSpotRepository(ctrl)
dateSpotRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
```

```go
// handler テストでの usecase mock 使用例
import usecasemock "github.com/daisuke-harada/date-courses-go/internal/usecase/mock"

ctrl := gomock.NewController(t)
defer ctrl.Finish()
mockPort := usecasemock.NewMockCreateDateSpotInputPort(ctrl)
mockPort.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&usecase.CreateDateSpotOutput{DateSpotID: 1}, nil)
```

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
- `mockPort.AssertNotCalled(t, "MethodName")` → 呼ばれないことの検証（testify/mock 使用時）

## 環境変数

`.envrc`（direnv）で管理。`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME` が必要。
