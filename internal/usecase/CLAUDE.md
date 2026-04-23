# internal/usecase/

## 実装方針

- `XxxInputPort` interface・`XxxInput`・`XxxOutput`・`XxxInteractor` を 1 ファイルにまとめる
- `openapi` パッケージへの依存禁止。handler 層で変換する
- `Output` はポインタ返し（`*XxxOutput`）
- エラーは必ず `apperror` パッケージ経由で返す
- `mock/` サブディレクトリに `package usecasemock` としてモックを自動生成する（`make gen`）

## コーディングルール

- `openapi` パッケージに依存しない。専用の `Input`/`Output` 型を定義する
- `Output` 型はポインタ（`*XxxOutput`）で返す
- エラーは `apperror` パッケージを使って返す（`apperror.NotFound()`, `apperror.InternalServerError(err)` など）
- repository の `SearchParams` 型に変換して渡す

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

## バリデーション（XxxInput の責務）

- **バリデーションは `XxxInput.Validate()` メソッドに実装する**
- `XxxInteractor.Execute` の冒頭で `input.Validate()` を呼ぶ
- 複数フィールドをまとめて検証し、エラーを `[]string` に収集してから `apperror.UnprocessableEntity(errs...)` で一括返却する
- handler 層にはバリデーションロジックを持たせない

```go
type XxxInput struct {
    Name  string
    Email string
}

func (i *XxxInput) Validate() error {
    var errs []string

    if strings.TrimSpace(i.Name) == "" {
        errs = append(errs, "名前を入力してください")
    } else if len(i.Name) > 50 {
        errs = append(errs, "名前は50文字以内で入力してください")
    }

    if len(errs) > 0 {
        return apperror.UnprocessableEntity(errs...)
    }
    return nil
}

func (i *XxxInteractor) Execute(ctx context.Context, input XxxInput) (*XxxOutput, error) {
    if err := input.Validate(); err != nil {
        return nil, err
    }
    // 以降はビジネスロジック
    ...
}
```
