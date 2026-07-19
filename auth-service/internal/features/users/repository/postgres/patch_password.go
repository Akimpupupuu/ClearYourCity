package users_repository

import (
	"context"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
)

func (r *usersRepository) PatchPassword(ctx context.Context, user core_domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	UPDATE auth_service.users
	SET
		hashed_password = $1,
		version = version+1
	WHERE id = $2 AND version = $3
	`

	cmdTag, err := r.pool.Exec(ctx, query, user.HashedPassword, user.ID, user.Version)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id = '%d' currently accessed: %w", user.ID, core_errors.ErrConflict)
	}

	return nil
}
