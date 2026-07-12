package sessions_service

import (
	"context"
	"fmt"
)

func (s *SessionsService) RevokeSession(ctx context.Context, tokenID string) error {
	err := s.sessionsRepository.RevokeSession(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("revoke session in repository: %w", err)
	}

	return nil
}
