package model

import "time"

// CourseAuthority はデートコースの公開設定です。
type CourseAuthority string

const (
	CourseAuthorityPublic  CourseAuthority = "公開"
	CourseAuthorityPrivate CourseAuthority = "非公開"
)

type Course struct {
	ID          uint            `gorm:"primaryKey;autoIncrement"`
	UserID      uint            `gorm:"not null;index"`
	TravelMode  string          `gorm:"not null"`
	Authority   CourseAuthority `gorm:"not null"`
	CreatedAt   time.Time       `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"not null;autoUpdateTime"`
	User        *User           `gorm:"foreignKey:UserID"`
	DuringSpots []*DuringSpot   `gorm:"foreignKey:CourseID"`
}
