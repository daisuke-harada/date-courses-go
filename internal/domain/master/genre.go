package master

// Genre は Rails の ActiveHash Genre と同等のマスタデータです。
type Genre struct {
	ID   int
	Name string
}

var genres = []Genre{
	{1, "ショッピングモール"}, {2, "飲食店"}, {3, "カフェ"},
	{4, "アウトドア"}, {5, "遊園地"}, {6, "水族館"},
	{7, "寿司"}, {8, "居酒屋"}, {9, "焼肉"},
	{10, "バーベキュー"}, {11, "ランドマーク"}, {12, "公園"},
}

// GenreNameByID は genre_id から名称を返します。存在しない ID は "" を返します。
func GenreNameByID(id int) string {
	for _, g := range genres {
		if g.ID == id {
			return g.Name
		}
	}
	return ""
}
