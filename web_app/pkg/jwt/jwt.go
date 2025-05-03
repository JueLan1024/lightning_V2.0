package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var ErrorInvalidToken = errors.New("invalid token")

var mySecret = []byte("不能说的")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"user_name"`
	jwt.StandardClaims
}

// genToken 生成JWT
func genToken(user_id int64, username string, time_duration time.Duration, issuer string) (string, error) {
	// 创建一个MyClaims
	claims := &MyClaims{
		user_id,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time_duration).Unix(),
			Issuer:    issuer,
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// GenAccessToken 生成JWT AccessToken
func GenAccessToken(user_id int64, username string, access_time_duration time.Duration) (string, error) {
	return genToken(user_id, username, access_time_duration, "access_token_issuer")
}

// GenRefreshToken 生成JWT RefreshToken
func GenRefreshToken(user_id int64, username string, refresh_time_duration time.Duration) (string, error) {
	return genToken(user_id, username, refresh_time_duration, "refresh_token_issuer")
}

// ParseToken解析tokenString
func ParseToken(tokenString string) (*MyClaims, error) {
	mc := new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(t *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return mc, nil
	}
	return nil, ErrorInvalidToken
}
