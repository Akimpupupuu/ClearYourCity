package sessions_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *SessionsRepository) GetSession(ctx context.Context, oldHashedToken string) (core_domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	SELECT id, user_id, refresh_token_hash, is_revoked, created_at, expires_at
	FROM auth_service.sessions
	WHERE refresh_token_hash = $1;
	`

	row := r.pool.QueryRow(ctx, query, oldHashedToken)

	var sessionModel SessionModel
	if err := row.Scan(
		&sessionModel.ID,
		&sessionModel.UserID,
		&sessionModel.RefreshTokenHash,
		&sessionModel.IsRevoked,
		&sessionModel.CreatedAt,
		&sessionModel.ExpiresAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.Session{}, fmt.Errorf("session with hash = '%s': %w", oldHashedToken, core_errors.ErrNotFound)
		}

		return core_domain.Session{}, fmt.Errorf("scan session: %w", err)
	}

	sessionDomain := DomainFromModel(sessionModel)
	return sessionDomain, nil
}
