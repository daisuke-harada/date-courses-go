package gemini

import "fmt"

// BuildDateSpotsPrompt はデートスポット検索用プロンプトを組み立てます。
// Gemini に対して JSON 形式での返答を指示します。
func BuildDateSpotsPrompt(prefectureName, genreName string, count int) string {
	return fmt.Sprintf(`%sの%sでデートにおすすめのスポットを%d件教えてください。
以下のJSON配列のみを返してください。説明文や追加テキストは不要です。

[
  {
    "name": "スポット名",
    "city_name": "市区町村名",
    "description": "簡単な説明（50文字以内）"
  }
]

要件:
- 実在する場所のみ
- name と city_name は必須
- JSON以外のテキストは一切含めない`, prefectureName, genreName, count)
}
