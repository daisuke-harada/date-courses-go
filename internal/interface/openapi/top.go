package openapi

import (
	"github.com/daisuke-harada/date-courses-go/internal/domain/master"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/samber/lo"
)

// NewTopResponse は master データと date spots から generated TopResponseData を組み立てます。
func NewTopResponse(dateSpots []*model.DateSpot) TopResponseData {
	return TopResponseData{
		DateSpots:       NewDateSpotSummaries(dateSpots),
		Areas:           newAreasResponse(),
		Genres:          newGenresResponse(),
		MainGenres:      newMainGenresResponse(),
		MainPrefectures: newMainPrefecturesResponse(),
	}
}

func newAreasResponse() []AreaData {
	areas := master.Areas()
	return lo.Map(areas, func(a master.Area, _ int) AreaData {
		return AreaData{Id: a.ID, Name: a.Name}
	})
}

func newGenresResponse() []GenreData {
	genres := master.Genres()
	return lo.Map(genres, func(g master.Genre, _ int) GenreData {
		return GenreData{Id: g.ID, Name: g.Name}
	})
}

func newMainGenresResponse() []GenreData {
	genres := master.MainGenres()
	return lo.Map(genres, func(g master.Genre, _ int) GenreData {
		return GenreData{Id: g.ID, Name: g.Name}
	})
}

func newMainPrefecturesResponse() []PrefectureData {
	prefectures := master.MainPrefectures()
	return lo.Map(prefectures, func(p master.Prefecture, _ int) PrefectureData {
		return PrefectureData{Id: p.ID, Name: p.Name, AreaId: p.AreaID}
	})
}
