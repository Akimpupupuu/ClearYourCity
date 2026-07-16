package sessions_jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type claimsKey struct{}

func NewUserClaims(userID int, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generate token uuid: %w", err)
	}

	return &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func ToContext(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

func FromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*UserClaims)
	return claims, ok
}
