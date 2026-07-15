package sessions_repository

import (
	"context"
	"fmt"
)

func (r *SessionsRepository) RevokeSession(ctx context.Context, hashedRefreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	UPDATE auth_service.sessions
	SET is_revoked = true
	WHERE refresh_token_hash = $1;
	`

	_, err := r.pool.Exec(ctx, query, hashedRefreshToken)
	if err != nil {
		return fmt.Errorf("execute revoke session: %w", err)
	}

	return nil
}
