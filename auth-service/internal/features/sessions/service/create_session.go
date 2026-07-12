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

	sum := sha256.Sum256([]byte(refreshToken))
	hashedRefreshToken := hex.EncodeToString(sum[:])

	sessionDomain := core_domain.NewSession(
		refreshClaims.ID,
		refreshClaims.UserID,
		hashedRefreshToken,
		false,
		refreshClaims.IssuedAt.Time,
		refreshClaims.ExpiresAt.Time,
	)

	if err := s.sessionsRepository.CreateSession(ctx, sessionDomain); err != nil {
		return SessionServiceResponse{}, fmt.Errorf("create session in repository: %w", err)
	}

	return SessionServiceResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
	}, nil
}
