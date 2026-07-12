package sessions_service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
)

func (s *SessionsService) CreateSession(ctx context.Context, userID int) (SessionServiceResponse, error) {
	accessToken, accessClaims, err := s.tokenGenerator.GenerateToken(userID, time.Minute*15)
	if err != nil {
		return SessionServiceResponse{}, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, refreshClaims, err := s.tokenGenerator.GenerateToken(userID, time.Hour*24)
	if err != nil {
		return SessionServiceResponse{}, fmt.Errorf("create refresh token: %w", err)
	}

	hash := sha256.New()
	hash.Write([]byte(refreshToken))
	hashedRefreshToken := hex.EncodeToString(hash.Sum(nil))

	sessionDomain := core_domain.NewSession(
		refreshClaims.ID,
		refreshClaims.UserID,
		hashedRefreshToken,
		false,
		refreshClaims.IssuedAt.Time,
		refreshClaims.ExpiresAt.Time,
	)

	session, err := s.sessionsRepository.CreateSession(ctx, sessionDomain)
	if err != nil {
		return SessionServiceResponse{}, fmt.Errorf("create session in repository: %w", err)
	}

	return SessionServiceResponse{
		session,
		accessToken,
		refreshToken,
		accessClaims.ExpiresAt.Time,
		refreshClaims.ExpiresAt.Time,
	}, nil
}
