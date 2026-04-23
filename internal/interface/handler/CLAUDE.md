# internal/interface/handler/

## 実装方針

- `XxxHandler` struct に `InputPort` フィールドを持つ
- フォームや JSON のパースはここで行い、usecase の `Input` 型に変換して渡す
- エラーは `return err` のみ（`ctx.JSON` でエラーを直接返さない）
- テストは同パッケージに `handler_test` パッケージで記述する

## コーディングルール

- `openapi` パッケージの型（`openapi.XxxParams`, `openapi.XxxJSONRequestBody` など）を受け取り、usecase の `Input` 型に変換して渡す
- エラーは `return err` のみ。`ctx.JSON` でエラーを直接返さない（`CustomHTTPErrorHandler` が処理する）
- usecase が返した `error` はそのまま `return err` する
- `XxxHandler` の `Xxx` メソッドは HTTP メソッドに対応する

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

## 型変換ルール

- handler はフォーム値の**型変換のみ**を行う。ビジネスルールのバリデーションは行わない
- `strconv.Atoi` などの型変換に失敗した場合は `apperror.BadRequest(...)` を返す
- `ctx.Bind(&req)` を使って struct にバインドする場合、型変換エラーは `return err` で返す

```go
// ✅ Good: 数値パース失敗は handler で BadRequest
id, err := strconv.Atoi(ctx.FormValue("user_id"))
if err != nil {
    return apperror.BadRequest("user_id は数値で指定してください")
}

// ✅ Good: ctx.Bind で型変換（バリデーションは usecase.Input.Validate() が担う）
if err := ctx.Bind(&req); err != nil {
    return err
}
input := usecase.XxxInput{Name: req.Name}
output, err := h.InputPort.Execute(ctx.Request().Context(), input)
```
