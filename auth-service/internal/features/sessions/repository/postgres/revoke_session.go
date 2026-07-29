package sessions_repository

import (
	"context"
	"fmt"

	core_postgres_transaction "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/postgres/transaction"
)

func (r *sessionsRepository) RevokeSession(ctx context.Context, hashedRefreshToken string) error {
	query := `
	UPDATE auth_service.sessions
	SET is_revoked = true
	WHERE refresh_token_hash = $1;
	`

	if tx, ok := core_postgres_transaction.GetTxFromContext(ctx); ok {
		_, err := tx.Exec(ctx, query, hashedRefreshToken)
		if err != nil {
			return fmt.Errorf("execute revoke session: %w", err)
		}

		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	_, err := r.pool.Exec(ctx, query, hashedRefreshToken)
	if err != nil {
		return fmt.Errorf("execute revoke session: %w", err)
	}

	return nil
}
