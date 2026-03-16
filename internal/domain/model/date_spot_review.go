package model

import "time"

type DateSpotReview struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	Rate       *float64
	Content    *string
	UserID     uint      `gorm:"not null;index"`
	DateSpotID uint      `gorm:"not null;index"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"not null;autoUpdateTime"`
	User       *User     `gorm:"foreignKey:UserID"`
	DateSpot   *DateSpot `gorm:"foreignKey:DateSpotID"`
}
