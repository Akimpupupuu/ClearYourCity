package core_jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type TokenGenerator struct {
	config Config
}

func NewTokenGenerator(config Config) *TokenGenerator {
	return &TokenGenerator{
		config: config,
	}
}

func (g *TokenGenerator) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(g.config.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	return claims, nil
}
