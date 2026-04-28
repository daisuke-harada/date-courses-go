package master

import "github.com/samber/lo"

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

// mainGenreIDs は Rails の Genre.majors に対応する ID スライスです。
var mainGenreIDs = []int{1, 2, 3, 4, 5, 6}

// GenreNameByID は genre_id から名称を返します。存在しない ID は "" を返します。
func GenreNameByID(id int) string {
	if g, ok := lo.Find(genres, func(g Genre) bool { return g.ID == id }); ok {
		return g.Name
	}
	return ""
}

// Genres returns all genre master data
func Genres() []Genre {
	return genres
}

func MainGenres() []Genre {
	return lo.Filter(genres, func(g Genre, _ int) bool {
		return lo.Contains(mainGenreIDs, g.ID)
	})
}
