package sessions_repository

import (
	"context"
	"fmt"

	core_postgres_transaction "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/postgres/transaction"
)

func (r *sessionsRepository) RevokeSessions(ctx context.Context, userID int) error {
	query := `
	UPDATE auth_service.sessions
	SET is_revoked = true
	WHERE user_id = $1;
	`

	if tx, ok := core_postgres_transaction.GetTxFromContext(ctx); ok {
		_, err := tx.Exec(ctx, query, userID)
		if err != nil {
			return fmt.Errorf("revoke sessions in repository: %w", err)
		}

		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("revoke sessions in repository: %w", err)
	}

	return nil
}
