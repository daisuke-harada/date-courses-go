# internal/domain/repository/

## 実装方針

- interface のみ定義。実装は `internal/infrastructure/persistence/` に置く
- メソッドは実際に使うものだけ定義する（YAGNI）
- 検索条件は `XxxSearchParams` 型にまとめる
- `mock/` サブディレクトリに `package repomock` としてモックを自動生成する（`make gen`）

## コーディングルール

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
```
