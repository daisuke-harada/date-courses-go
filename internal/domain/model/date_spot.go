package model

import "time"

type DateSpot struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	GenreID      *int   `gorm:"index"`
	PrefectureID *int   `gorm:"index"`
	Name         string `gorm:"not null"`
	CityName     string `gorm:"not null"`
	Image        *string
	Latitude     *float64
	Longitude    *float64
	OpeningTime  *time.Time
	ClosingTime  *time.Time
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`

	// DB集計フィールド (SELECT時のみ使用、マイグレーション対象外)
	// gorm:"->" だと Find() でカラムがマッピングされないケースがあるため column タグのみ指定する
	AverageRate       float64 `gorm:"column:average_rate;<-:false"`
	ReviewTotalNumber int     `gorm:"column:review_total_number;<-:false"`
}
