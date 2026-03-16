package model

import "time"

type Course struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null;index"`
	TravelMode string    `gorm:"not null"`
	Authority  string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"not null;autoUpdateTime"`
	User       *User     `gorm:"foreignKey:UserID"`
}
