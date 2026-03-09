package domain

import "time"

type Relationship struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	FollowID  uint      `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
	User      *User     `gorm:"foreignKey:UserID"`
	Follow    *User     `gorm:"foreignKey:FollowID"`
}
