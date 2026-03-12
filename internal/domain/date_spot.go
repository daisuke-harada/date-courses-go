package domain

import "time"

type DateSpot struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	GenreID     *int   `gorm:"index"`
	Name        string `gorm:"not null"`
	Image       *string
	OpeningTime *time.Time
	ClosingTime *time.Time
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}
