package model

import "time"

type Address struct {
	ID           uint `gorm:"primaryKey;autoIncrement"`
	PrefectureID *int
	DateSpotID   *int   `gorm:"index"`
	CityName     string `gorm:"not null"`
	Latitude     *float64
	Longitude    *float64
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
	DateSpot     *DateSpot `gorm:"foreignKey:DateSpotID"`
}
