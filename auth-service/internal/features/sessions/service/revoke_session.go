package sessions_service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func (s *SessionsService) RevokeSession(ctx context.Context, refreshToken string) error {
	sum := sha256.Sum256([]byte(refreshToken))
	hashedRefreshToken := hex.EncodeToString(sum[:])

	err := s.sessionsRepository.RevokeSession(ctx, hashedRefreshToken)
	if err != nil {
		return fmt.Errorf("revoke session in repository: %w", err)
	}

	return nil
}
