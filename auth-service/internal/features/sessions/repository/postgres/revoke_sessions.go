package sessions_repository

import (
	"context"
	"fmt"
)

func (r *sessionsRepository) RevokeSessions(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	UPDATE auth_service.sessions
	SET is_revoked = true
	WHERE user_id = $1;
	`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("revoke sessions in repository: %w", err)
	}

	return nil
}
