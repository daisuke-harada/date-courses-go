---
name: design-then-code
description: 設計を上位モデル（Opus）で行い、実装だけを下位モデル（Sonnet）の go-implementer サブエージェントに任せる開発ワークフロー。新しい API・機能・リファクタの実装依頼を受けたときに使う。「設計してから実装して」「設計はOpusで実装はSonnetで」「新しいAPIを追加して」「この機能を実装して」と言われたら使う。
---

# design-then-code — 設計は上位モデル、実装は下位モデル

このスキルは、実装作業を **「設計（あなた＝上位モデル）」** と **「コーディング（`go-implementer` サブエージェント＝Sonnet）」** に分離するためのものです。
あなたは **コードを自分で書かない**。設計し、指示書を書き、実装をレビューする役に徹します。

## なぜ分けるのか

- 設計判断（レイヤー分割・責務・API 契約・エラー設計）は上位モデルの精度が要る。
- 規約に沿ったコードを書く作業は、良い指示書さえあれば下位モデルで十分にこなせる。
- 分離することで、上位モデルのコンテキストを「設計とレビュー」に集中して使える。

## 手順

### Step 1: 調査（あなたが行う）

実装対象を理解するために、必要な範囲を自分で読む。

- 関連する既存コード（同じレイヤーの類似実装、呼び出し元、DB スキーマ）
- `api/` の OpenAPI スキーマ
- ルートおよび対象ディレクトリの `CLAUDE.md`

広く探す必要があるときは `Explore` エージェントを使ってよい。**ここで設計に必要な事実を固めきる** — 指示書の精度がそのまま実装品質になります。

### Step 2: 設計する（あなたが行う）

次を決める。曖昧さを残さない。

- **API 契約**: パス・メソッド・request / response のフィールドと型・ステータスコード（変更がある場合）
- **レイヤー構成**: どのファイルに何を置くか（model / repository / usecase / handler / DI）
- **データフロー**: handler → usecase → repository の入出力の型
- **エラー設計**: どの条件で `apperror` のどの関数を返すか
- **テストケース**: 正常系・異常系を具体名で列挙（`success`, `error_not_found` …）

判断が分かれる設計上の選択（テーブル追加の要否、既存 API の破壊的変更など）が出たら、**実装に入る前にユーザーへ確認する**。実装が終わってからの手戻りが一番高くつきます。

### Step 3: 実装指示書を書く

scratchpad ディレクトリに `design-<機能名>.md` として書き出す。テンプレートは後述。

### Step 4: `go-implementer` に実装させる

`Agent` ツールで実装を委譲する。

```
subagent_type: "go-implementer"
run_in_background: false        // 次のレビューが結果に依存するため
prompt: 指示書の全文（ファイルパスだけでなく、本文をそのまま渡す）
```

ポイント:

- **指示書の中身をプロンプトに丸ごと入れる。** サブエージェントはコンテキストを引き継がないため、「さっき話した件」は通じません。
- 独立した機能が複数あるなら、**並列に複数体起動してよい**（同じファイルを触るタスク同士は必ず直列にする）。
- 依存関係のあるタスク（issue-N1 → N2）は、branch chaining の順序どおり直列に実行する。

### Step 5: レビューする（あなたが行う）

サブエージェントの報告は**そのまま信用しない**。必ず自分で確認する。

1. `git diff` で実際の差分を読む。
2. 設計どおりか照合する — レイヤー違反、規約違反（生成物の直接編集、テストなし実装、日本語テスト名など）がないか。
3. 報告された検証結果を自分で再実行して裏を取る。

```sh
GOROOT=/Users/haradadaisuke/.goenv/versions/1.26.2 make lint
GOROOT=/Users/haradadaisuke/.goenv/versions/1.26.2 make test
```

4. 修正が必要なら、**`SendMessage` で同じサブエージェントに差し戻す**（新規 Agent 起動はコンテキストが失われるので避ける）。設計そのものの誤りだった場合は、指示書を直してから差し戻す。
5. 小さな修正（typo、import 1 行）は、差し戻すより自分で直したほうが速い。判断はあなたに任せます。

### Step 6: 報告

ユーザーには「設計の要点」「実装された内容」「lint / test の結果」を伝える。
サブエージェントの報告は**ユーザーには見えていない** ので、必要な内容は自分の言葉で伝えること。

commit / push / PR 作成は、ユーザーが明示的に求めた場合のみ行う。

---

## 実装指示書テンプレート

```markdown
# 実装指示書: <機能名>

## ゴール
<1〜3行。何ができるようになるか>

## 前提・背景
- 関連する既存実装: <path:line>
- 参照すべき CLAUDE.md: <path>
- 既存の類似実装（真似すべきお手本）: <path>

## OpenAPI 変更
（不要なら「なし」と明記する）
- ファイル: api/paths/xxx.yaml
- エンドポイント: POST /xxx
- request:
  | フィールド | 型 | 必須 | 備考 |
- response (200):
  | フィールド | 型 | 備考 |
- 認証: 必要 / 不要（必要なら security: [{ bearerAuth: [] }]）

## 変更対象ファイル
| ファイル | 変更内容 |
|---|---|
| internal/usecase/xxx_interactor.go | 新規。XxxInteractor.Execute を実装 |
| ... | ... |

## 実装詳細

### usecase
- InputPort / OutputPort のシグネチャ:
  ```go
  type XxxInputPort interface {
      Execute(ctx context.Context, input XxxInput) (*XxxOutput, error)
  }
  ```
- ビジネスルール:
  1. <条件> のとき <振る舞い>
- エラー:
  | 条件 | 返すエラー |
  |---|---|
  | 対象が存在しない | apperror.NotFound() |

### handler
- バリデーション: <ルール>
- 型変換: <どの変換関数を使うか>

### DI
- di/infrastructure.go に <何> を登録

## テストケース（この名前で書くこと）
### internal/usecase/xxx_interactor_test.go
- TestXxxInteractor_Execute
  - success
  - error_not_found
  - error_repository_failure

### internal/interface/handler/xxx_handler_test.go
- TestXxxHandler_Xxx
  - success
  - error_invalid_request
  - error_usecase_failure

## 完了条件
- 上記テストがすべて通る
- `make gen && go build ./... && make lint && make test` がすべてエラーゼロ
```

---

## 守ること

- **あなたはコードを書かない。** 実装は `go-implementer` に任せる。例外は、レビューで見つけた 1〜2 行の軽微な修正のみ。
- **曖昧な指示書を渡さない。** 「よしなに実装して」は失敗する。ファイル名・関数シグネチャ・テストケース名まで具体的に書く。
- **サブエージェントの「完了しました」を鵜呑みにしない。** 差分と lint / test 結果を自分で確認する。
- **設計上の疑問はユーザーに聞く。** 勝手に決めて実装させてから直すのが一番高くつく。
