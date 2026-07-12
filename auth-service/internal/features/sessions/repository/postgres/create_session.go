package sessions_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgxViolatesForeignKeyErrorCode = "23503"
)

func (r *SessionsRepository) CreateSession(ctx context.Context, session core_domain.Session) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	INSERT INTO auth_service.sessions (id, user_id, refresh_token_hash, is_revoked, created_at, expires_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.IsRevoked,
		session.CreatedAt,
		session.ExpiresAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgxViolatesForeignKeyErrorCode {
				return fmt.Errorf("%v: user with id = '%d': %w", err, session.UserID, core_errors.ErrNotFound)
			}
		}

		return fmt.Errorf("execute create session: %w", err)
	}

	return nil
}
