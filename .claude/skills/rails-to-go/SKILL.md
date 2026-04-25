# Rails → Go 移植スキル

このファイルは Claude エージェントが `date-courses-go` の Issue を自律実装する際に蓄積したテクニックを記録します。

---

## 学習済みテクニック

- [2026-04-25] open Issue リストには PR も含まれるため、`pull_request` キーの有無で「実際の Issue か PR か」を判別する必要がある。Issue のみを実装対象とする。
- [2026-04-25] Issue に既存の open PR が紐付いている場合（PR タイトルに `closes #N` など）は、重複実装を避けるためスキップし次の Issue に進む。
- [2026-04-25] Notion MCP が利用不可の場合、セッションサマリーはテキスト出力のみ行い、次回セッションで手動登録するよう注記する。
- [2026-04-25] `.claude/skills/rails-to-go/SKILL.md` が存在しない場合は `mcp__github__create_or_update_file` で新規作成する（sha 不要）。既存の場合は sha を取得してから更新する。
- [2026-04-25] main ブランチへの直接プッシュが保護されている場合、SKILL.md 更新も別ブランチ → PR 経由で行う必要がある。
