package service

import (
	"golang.org/x/crypto/bcrypt"
)

// AuthService はパスワードのハッシュ化・検証を行うドメインサービスです。
// Rails の has_secure_password (bcrypt) と同等の処理を提供します。
type AuthService interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) bool
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

// HashPassword はパスワードを bcrypt でハッシュ化します。
func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword はハッシュ化されたパスワードと平文パスワードを比較します。
func (s *authService) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
