package jwt

import (
	"errors"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func Encode(userID uint, secretKey string) (string, error) {
	return EncodeWithExpiry(userID, secretKey, time.Now().Add(24*time.Hour))
}

func EncodeWithExpiry(userID uint, secretKey string, expiry time.Time) (string, error) {
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secretKey))
}

func Decode(tokenStr string, secretKey string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, apperror.Unauthorized("トークンの有効期限が切れています。再度ログインしてください。")
		}
		return 0, apperror.Unauthorized("認証に失敗しました。")
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return 0, apperror.Unauthorized("認証に失敗しました。")
	}
	return c.UserID, nil
}
