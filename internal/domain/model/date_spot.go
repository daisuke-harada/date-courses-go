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

	// Gemini 統合フィールド
	// source: 'manual'（管理者手動登録）or 'ai'（Gemini 自動生成）
	Source string `gorm:"not null;default:manual"`
	// maps_url: Google Maps 検索 URL（AI スポットはこちらで詳細確認）
	MapsURL *string `gorm:"column:maps_url"`
	// link_status: 'active' / 'invalid' / 'unchecked'（週次バッチで更新）
	LinkStatus string `gorm:"not null;default:active"`
	// last_checked_at: リンク最終確認日時
	LastCheckedAt *time.Time
	// normalized_name: 重複チェック用（全角半角統一・スペース除去・小文字化）
	NormalizedName string

	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`

	// DB集計フィールド (SELECT時のみ使用、マイグレーション対象外)
	// gorm:"->" だと Find() でカラムがマッピングされないケースがあるため column タグのみ指定する
	AverageRate       float64 `gorm:"column:average_rate;<-:false"`
	ReviewTotalNumber int     `gorm:"column:review_total_number;<-:false"`
}
