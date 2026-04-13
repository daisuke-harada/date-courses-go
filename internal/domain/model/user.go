package model

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	Name           string    `gorm:"not null;uniqueIndex"`
	Email          string    `gorm:"not null;uniqueIndex"`
	Gender         string    `gorm:"not null"`
	Image          *string
	Admin          bool      `gorm:"not null;default:false"`
	PasswordDigest string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime"`
}
