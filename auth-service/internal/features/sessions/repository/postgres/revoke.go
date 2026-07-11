package sessions_repository

import (
	"context"
	"fmt"
)

func (r *SessionsRepository) RevokeSession(ctx context.Context, tokenID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	UPDATE auth_service.sessions
	SET is_revoked = true
	WHERE id = $1;
	`

	_, err := r.pool.Exec(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("exequte revoke session: %w", err)
	}

	return nil
}
