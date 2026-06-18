package model

import "time"

type DateSpotSource string
type DateSpotLinkStatus string

const (
	DateSpotSourceManual DateSpotSource = "manual"
	DateSpotSourceAI     DateSpotSource = "ai"

	DateSpotLinkStatusActive    DateSpotLinkStatus = "active"
	DateSpotLinkStatusInvalid   DateSpotLinkStatus = "invalid"
	DateSpotLinkStatusUnchecked DateSpotLinkStatus = "unchecked"
)

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
	Source         DateSpotSource     `gorm:"not null;default:manual"`
	MapsURL        *string            `gorm:"column:maps_url"`
	LinkStatus     DateSpotLinkStatus `gorm:"not null;default:active"`
	LastCheckedAt  *time.Time
	NormalizedName string
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`

	// DB集計フィールド (SELECT時のみ使用、マイグレーション対象外)
	// gorm:"->" だと Find() でカラムがマッピングされないケースがあるため column タグのみ指定する
	AverageRate       float64 `gorm:"column:average_rate;<-:false"`
	ReviewTotalNumber int     `gorm:"column:review_total_number;<-:false"`
}
