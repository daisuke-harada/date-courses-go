package model

import "time"

type Gender string

const (
	GenderMale   Gender = "男性"
	GenderFemale Gender = "女性"
)

type User struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	Name           string `gorm:"not null;uniqueIndex"`
	Email          string `gorm:"not null;uniqueIndex"`
	Gender         Gender `gorm:"not null"`
	Image          *string
	Admin          bool      `gorm:"not null;default:false"`
	PasswordDigest string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime"`
}
