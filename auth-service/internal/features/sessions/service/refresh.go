package sessions_service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
)

func (s *SessionsService) RefreshToken(ctx context.Context, oldRefreshToken string) (SessionServiceResponse, error) {
	hash := sha256.New()
	hash.Write([]byte(oldRefreshToken))
	hashedRefreshToken := hex.EncodeToString(hash.Sum(nil))

	session, err := s.sessionsRepository.GetSession(ctx, hashedRefreshToken)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return SessionServiceResponse{}, fmt.Errorf("get session from repository: %v: %w", err, core_errors.ErrUnauthorized)
		}
		return SessionServiceResponse{}, fmt.Errorf("get session from repository: %w", err)
	}

	if err := s.RevokeSession(ctx, session.ID); err != nil {
		return SessionServiceResponse{}, fmt.Errorf("revoke session: %w", err)
	}

	accessToken, accessClaims, err := s.tokenGenerator.GenerateToken(session.UserID, time.Minute*15)
	if err != nil {
		return SessionServiceResponse{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, refreshClaims, err := s.tokenGenerator.GenerateToken(session.UserID, time.Hour*24)
	if err != nil {
		return SessionServiceResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}

	hash.Reset()
	hash.Write([]byte(refreshToken))
	hashedRefreshToken = hex.EncodeToString(hash.Sum(nil))

	sessionDomain := core_domain.NewSession(
		refreshClaims.ID,
		refreshClaims.UserID,
		hashedRefreshToken,
		false,
		refreshClaims.IssuedAt.Time,
		refreshClaims.ExpiresAt.Time,
	)

	newSession, err := s.sessionsRepository.CreateSession(ctx, sessionDomain)
	if err != nil {
		return SessionServiceResponse{}, fmt.Errorf("create new session in repository: %w", err)
	}

	return SessionServiceResponse{
		newSession,
		accessToken,
		refreshToken,
		accessClaims.ExpiresAt.Time,
		refreshClaims.ExpiresAt.Time,
	}, nil
}
