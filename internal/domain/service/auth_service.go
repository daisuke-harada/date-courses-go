package service

import "golang.org/x/crypto/bcrypt"

// AuthService はパスワードのハッシュ化・検証を担うドメインサービスです。
// Rails の has_secure_password（bcrypt）相当の機能を提供します。
type AuthService interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) bool
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

// HashPassword はプレーンテキストのパスワードを bcrypt でハッシュ化します。
func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword はハッシュ済みパスワードとプレーンテキストを比較します。
// Rails の user.authenticate(password) に相当します。
func (s *authService) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
