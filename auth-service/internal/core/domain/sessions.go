package core_domain

import (
	"fmt"
	"time"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
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

func (s *Session) Validate() error {
	if s.IsRevoked {
		return fmt.Errorf("token reuse: %w", core_errors.ErrTokenReuse)
	}

	if s.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token expired: %w", core_errors.ErrUnauthorized)
	}

	return nil
}
