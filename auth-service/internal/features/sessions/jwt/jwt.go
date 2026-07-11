package sessions_jwt

import (
	"fmt"
	"time"

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

func (g *TokenGenerator) GenerateToken(userID int, duration time.Duration) (string, *UserClaims, error) {
	userClaims, err := NewUserClaims(userID, duration)
	if err != nil {
		return "", nil, fmt.Errorf("create user claims: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenStr, err := token.SignedString([]byte(g.config.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("sign token: %w", err)
	}

	return tokenStr, userClaims, nil
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
