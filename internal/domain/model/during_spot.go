package model

import "time"

type DuringSpot struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	CourseID   uint      `gorm:"not null;index"`
	DateSpotID uint      `gorm:"not null;index"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"not null;autoUpdateTime"`
	Course     *Course   `gorm:"foreignKey:CourseID"`
	DateSpot   *DateSpot `gorm:"foreignKey:DateSpotID"`
}
