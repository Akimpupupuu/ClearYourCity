package sessions_service

import (
	"context"
	"fmt"
)

func (s *SessionsService) RevokeSessions(ctx context.Context, userID int) error {
	err := s.sessionsRepository.RevokeSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("revoke sessions: %w", err)
	}

	return nil
}
