# internal/domain/service/

## 実装方針

- 複数の repository をまたぐドメインロジックをインターフェースで定義する
- `mock/` サブディレクトリに `package servicemock` としてモックを自動生成する（`make gen`）
