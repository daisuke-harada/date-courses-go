package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

// NewUser は新規ユーザーを生成します。
// passwordDigest にはハッシュ化済みのパスワードを渡してください。
// image は nil を許容します。
func NewUser(name string, email string, gender Gender, image *string, passwordDigest string) *User {
	return &User{
		Name:           name,
		Email:          email,
		Gender:         gender,
		Image:          image,
		PasswordDigest: passwordDigest,
	}
}

// ApplyUpdate はユーザーの更新可能なフィールドを上書きします。
// password が空文字の場合はパスワードを更新しません。
// image は nil の場合は更新しません。
func (u *User) ApplyUpdate(name string, email string, gender Gender, image *string, password string) error {
	u.Name = name
	u.Email = email
	u.Gender = gender
	if image != nil {
		u.Image = image
	}
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordDigest = string(hashed)
	}
	return nil
}

// UserWithRelations はユーザーと関連データをまとめた中間型です。
type UserWithRelations struct {
	User         *User
	FollowerIDs  []int
	FollowingIDs []int
	Courses      []*Course
	Reviews      []*DateSpotReview
}
