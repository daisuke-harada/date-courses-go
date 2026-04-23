# internal/domain/model/

## 実装方針

- GORM タグ付きの struct を定義する
- ビジネスロジックは持たない（値オブジェクトは型エイリアスで表現）
- 例: `type Gender string` + `const GenderMale Gender = "male"`
- 複数のモデルをまとめた集約型（例: `XxxWithRelations`）もここに定義する
- usecase や openapi など他のパッケージの型には依存しない

## 列挙型（enum）に関する注意

ドメインモデルに列挙型を追加する場合は以下を守ってください。

- 型は基本的に型エイリアスで定義する（例: `type Gender string`）
- 値は `const` で列挙する（例: `const ( GenderMale Gender = "male" ... )`）
- JSON エンコード/デコードや DB マイグレーションで特別な処理が必要な場合は、`MarshalJSON`/`UnmarshalJSON` を実装するか、変換ヘルパーを用意する
- OpenAPI 側に別の enum 型がある場合は、`internal/interface/openapi/` 側に変換関数（`NewXxx` / `ToModelXxx`）を実装して一元管理すること
- 型の追加・変更を行ったら、該当するテストと `make gen`（openapi/types 再生成や mock 再生成）が必要になる点に注意する
