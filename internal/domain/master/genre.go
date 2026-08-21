package master

import "github.com/samber/lo"

// Genre は Rails の ActiveHash Genre と同等のマスタデータです。
type Genre struct {
	ID   int
	Name string
}

// HotPepper グルメAPI のジャンルに 1:1 で対応させたデートグルメのジャンル。
var genres = []Genre{
	{1, "居酒屋"}, {2, "ダイニングバー・バル"}, {3, "カフェ・スイーツ"},
	{4, "和食"}, {5, "洋食"}, {6, "イタリアン・フレンチ"},
	{7, "中華"}, {8, "焼肉・ホルモン"}, {9, "ラーメン"},
	{10, "アジア・エスニック料理"}, {11, "韓国料理"}, {12, "バー・カクテル"},
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
