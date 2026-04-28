package master

import "github.com/samber/lo"

// Prefecture は Rails の ActiveHash Prefecture と同等のマスタデータです。
type Prefecture struct {
	ID     int
	Name   string
	AreaID int
}

var mainPrefectureIDs = []int{13, 27, 40, 14, 23, 26}

var prefectures = []Prefecture{
	{1, "北海道", 1}, {2, "青森県", 1}, {3, "岩手県", 1}, {4, "宮城県", 1}, {5, "秋田県", 1},
	{6, "山形県", 1}, {7, "福島県", 1}, {8, "茨城県", 2}, {9, "栃木県", 2}, {10, "群馬県", 2},
	{11, "埼玉県", 2}, {12, "千葉県", 2}, {13, "東京都", 2}, {14, "神奈川県", 2}, {15, "新潟県", 3},
	{16, "富山県", 3}, {17, "石川県", 3}, {18, "福井県", 3}, {19, "山梨県", 3}, {20, "長野県", 3},
	{21, "岐阜県", 3}, {22, "静岡県", 3}, {23, "愛知県", 3}, {24, "三重県", 3}, {25, "滋賀県", 4},
	{26, "京都府", 4}, {27, "大阪府", 4}, {28, "兵庫県", 4}, {29, "奈良県", 4}, {30, "和歌山県", 4},
	{31, "鳥取県", 5}, {32, "島根県", 5}, {33, "岡山県", 5}, {34, "広島県", 5}, {35, "山口県", 5},
	{36, "徳島県", 5}, {37, "香川県", 5}, {38, "愛媛県", 5}, {39, "高知県", 5}, {40, "福岡県", 6},
	{41, "佐賀県", 6}, {42, "長崎県", 6}, {43, "熊本県", 6}, {44, "大分県", 6}, {45, "宮崎県", 6},
	{46, "鹿児島県", 6}, {47, "沖縄県", 6},
}

// PrefectureNameByID は prefecture_id から名称を返します。存在しない ID は "" を返します。
func PrefectureNameByID(id int) string {
	for _, p := range prefectures {
		if p.ID == id {
			return p.Name
		}
	}
	return ""
}

func PrefectureByID(id int) *Prefecture {
	for _, p := range prefectures {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

// Prefectures returns all prefecture master data
func Prefectures() []Prefecture {
	return prefectures
}

func MainPrefectures() []Prefecture {
	return lo.Filter(prefectures, func(p Prefecture, _ int) bool {
		return lo.Contains(mainPrefectureIDs, p.ID)
	})
}
