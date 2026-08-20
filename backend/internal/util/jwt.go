package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 签发方，登录签发与校验必须一致。
const jwtIssuer = "gbevent"

// Claims 自定义 JWT 载荷。
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT。
func GenerateToken(secret string, expire time.Duration, userID uint64, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// ParseToken 解析并校验 JWT：签名、过期时间、签发方均须合法，否则返回错误。
func ParseToken(secret, tokenString string) (*Claims, error) {
	if secret == "" || tokenString == "" {
		return nil, errors.New("empty secret or token")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		// 仅接受 HMAC，杜绝 alg=none 及算法混淆攻击。
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(secret), nil
	}, jwt.WithIssuer(jwtIssuer))
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
