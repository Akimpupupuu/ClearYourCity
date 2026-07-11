package core_domain

import (
	"time"
)

type Session struct {
	ID           string
	UserID       int
	RefreshToken string
	IsRevoked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

func NewSession(
	id string,
	userID int,
	refreshToken string,
	isRevoked bool,
	createdAt time.Time,
	expiresAt time.Time,
) Session {
	return Session{
		ID:           id,
		UserID:       userID,
		RefreshToken: refreshToken,
		IsRevoked:    isRevoked,
		CreatedAt:    createdAt,
		ExpiresAt:    expiresAt,
	}
}
