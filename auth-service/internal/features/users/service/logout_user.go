package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) LogoutUser(ctx context.Context, refreshToken string) error {
	if err := s.sessionsService.RevokeSession(ctx, refreshToken); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
