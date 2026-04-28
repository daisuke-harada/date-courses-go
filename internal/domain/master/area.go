package master

// Area は地域マスタデータ（例: 北海道・東北、関東 など）
type Area struct {
	ID   int
	Name string
}

var areas = []Area{
	{1, "北海道・東北"},
	{2, "関東"},
	{3, "中部"},
	{4, "関西"},
	{5, "中国・四国"},
	{6, "九州・沖縄"},
}

// Areas returns list of area master data
func Areas() []Area {
	return areas
}
