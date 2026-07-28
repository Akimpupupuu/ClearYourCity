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
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
)

func (s *sessionsService) RefreshToken(ctx context.Context, oldRefreshToken string) (SessionServiceResponse, error) {
	sum := sha256.Sum256([]byte(oldRefreshToken))
	hashedRefreshToken := hex.EncodeToString(sum[:])

	session, err := s.sessionsRepository.GetSession(ctx, hashedRefreshToken)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return SessionServiceResponse{}, fmt.Errorf("get session from repository: %v: %w", err, core_errors.ErrUnauthorized)
		}
		return SessionServiceResponse{}, fmt.Errorf("get session from repository: %w", err)
	}

	if err := session.Validate(); err != nil {
		if errors.Is(err, core_errors.ErrTokenReuse) {
			_ = s.sessionsRepository.RevokeSessions(ctx, session.UserID)
			return SessionServiceResponse{}, fmt.Errorf("token reuse detected: %w", err)
		}

		return SessionServiceResponse{}, fmt.Errorf("validate session: %w", err)
	}

	var (
		accessToken   string
		refreshToken  string
		accessClaims  *sessions_jwt.UserClaims
		refreshClaims *sessions_jwt.UserClaims
	)

	if err = s.transactionManager.WithTransaction(ctx, func(ctx context.Context) error {
		var err error

		if err = s.sessionsRepository.RevokeSession(ctx, hashedRefreshToken); err != nil {
			return fmt.Errorf("revoke session: %w", err)
		}

		accessToken, accessClaims, err = s.tokenGenerator.GenerateToken(session.UserID, time.Minute*15)
		if err != nil {
			return fmt.Errorf("generate access token: %w", err)
		}

		refreshToken, refreshClaims, err = s.tokenGenerator.GenerateToken(session.UserID, time.Hour*24)
		if err != nil {
			return fmt.Errorf("generate refresh token: %w", err)
		}

		sum = sha256.Sum256([]byte(refreshToken))
		hashedRefreshToken = hex.EncodeToString(sum[:])

		sessionDomain := core_domain.NewSession(
			refreshClaims.ID,
			refreshClaims.UserID,
			hashedRefreshToken,
			false,
			refreshClaims.IssuedAt.Time,
			refreshClaims.ExpiresAt.Time,
		)

		if err = s.sessionsRepository.CreateSession(ctx, sessionDomain); err != nil {
			return fmt.Errorf("create new session in repository: %w", err)
		}

		return nil
	}); err != nil {
		return SessionServiceResponse{}, fmt.Errorf("transaction error: %w", err)
	}

	return SessionServiceResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
	}, nil
}
