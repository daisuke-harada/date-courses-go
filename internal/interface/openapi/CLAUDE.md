# internal/interface/openapi/

## 実装方針

- `api_server.gen.go`, `api_types.gen.go` は **編集禁止**（`make gen` で再生成）
- レスポンス変換用の手書きファイル（`signup.go`, `login.go` など）はここに置く
- `Gender` 型など openapi 型への変換は `NewXxx()` 関数を定義する

## レスポンス整形ルール

- usecase の `Output` → openapi の型への変換は **このディレクトリに手書きの変換関数を定義する**
- 関数名は `NewXxxResponse` に統一する
- handler 側では変換関数を呼び出すだけにし、整形ロジックを持たせない
- 変換中にエラーが起きた場合は `apperror.InternalServerError(err)` で返す

```go
// internal/interface/openapi/date_spot.go
func NewDateSpotResponse(ds *model.DateSpot) (DateSpotResponseBody, error) {
    return DateSpotResponseBody{
        Id:   int(ds.ID),
        Name: ds.Name,
        // ...
    }, nil
}

// handler 側
output, err := h.InputPort.Execute(ctx.Request().Context(), input)
if err != nil {
    return err
}
resp, err := openapi.NewDateSpotResponse(output.DateSpot)
if err != nil {
    return apperror.InternalServerError(err)
}
return ctx.JSON(http.StatusOK, resp)
```
