package jwtprocessing

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

const (
	TOKENEXPIRES = 240 * time.Hour
)

func GenerateToken(login, key string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXPIRES)),
		},
		Login: login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(key))
}

// ParseToken - parse token
func ParseToken(tokenString, key string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if claims.Valid() != nil {
		return "", claims.Valid()
	}
	return claims.Login, err
}

// Valid - check if token is valid
func (c *Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
