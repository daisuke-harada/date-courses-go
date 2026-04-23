# internal/infrastructure/persistence/

## 実装方針

- `internal/domain/repository/` インターフェースの GORM 実装のみ置く
- GORM の WHERE 句はここだけに閉じ込める
- エラーは `apperror.InternalServerError(err)` か `apperror.NotFound()` に変換して返す

```go
func (r *xxxRepository) Search(ctx context.Context, params repository.XxxSearchParams) ([]*model.Xxx, error) {
    db := r.db.WithContext(ctx).Model(&model.Xxx{})
    if params.Name != nil {
        db = db.Where("name LIKE ?", "%"+*params.Name+"%")
    }
    var result []*model.Xxx
    if err := db.Find(&result).Error; err != nil {
        return nil, apperror.InternalServerError(err)
    }
    return result, nil
}
```
