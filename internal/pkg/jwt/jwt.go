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
	// jwt.ParseWithClaims は次の処理を行います：
	// 1) トークンを 3 部分に分割し、ヘッダー／ペイロードをデコードして構造体に Unmarshal する
	// 2) ヘッダーの "alg" から署名方式を決定し、ここで渡す keyfunc を呼んで検証用キーを取得する
	// 3) 取得したキーで署名を検証する（署名不正ならエラー）
	// 4) Claims.Valid() を呼んで RegisteredClaims の検証（exp / nbf 等）を行う
	//    - 期限切れの場合は jwt.ErrTokenExpired が返される
	// つまり署名検証と有効期限チェックはこの呼び出し内で行われます。
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (any, error) {
		// keyfunc: ライブラリが署名検証のために呼び出します。
		// ここでは署名方式が HMAC 系であることを確認し、検証に使う対称鍵を返します。
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// 実際の検証はライブラリ側で行われます。ここではキーを返すだけです。
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, apperror.Unauthorized("トークンの有効期限が切れています。再度ログインしてください。")
		}
		return 0, apperror.Unauthorized("認証に失敗しました。")
	}

	// パース済みトークンからクレームを取り出し、妥当性を確認する
	// (jwt.ParseWithClaims で署名検証と期限チェックは既に行われている)
	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return 0, apperror.Unauthorized("認証に失敗しました。")
	}
	return c.UserID, nil
}
