package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
			Issuer:    "gbevent",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// ParseToken 解析并校验 JWT。
func ParseToken(secret, tokenString string) (*Claims, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	claims := token.Claims.(*Claims)
	return claims, nil
}
