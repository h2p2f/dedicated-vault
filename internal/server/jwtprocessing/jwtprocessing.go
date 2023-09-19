package jwtprocessing

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	Login    string
	ClientID string
}

const (
	TOKENEXPIRES = 240 * time.Hour
)

func GenerateToken(login, agent, key string) (string, error) {
	fmt.Println("generating token", login, agent, key)
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXPIRES)),
		},
		Login:    login,
		ClientID: agent,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(key))
}

// ParseToken - parse token
func ParseToken(tokenString, key string) (string, string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if claims.Valid() != nil {
		return "", "", claims.Valid()
	}
	return claims.Login, claims.ClientID, err
}

// Valid - check if token is valid
func (c *Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
