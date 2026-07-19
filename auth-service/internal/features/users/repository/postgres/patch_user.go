package users_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *usersRepository) PatchUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	UPDATE auth_service.users 
	SET
		full_name = $1,
		email = $2,
		version = version+1
	WHERE id = $3 AND version = $4
	RETURNING id, version, full_name, email, hashed_password, created_at
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.Email, user.ID, user.Version)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.Email,
		&userModel.HashedPassword,
		&userModel.CreatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgxViolatesUniqueErrorCode {
				return core_domain.User{}, fmt.Errorf("%v: user with email = '%s': %w", err, user.Email, core_errors.ErrConflict)
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf("user with id = '%d' concurrently accessed: %w", user.ID, core_errors.ErrConflict)
		}

		return core_domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	userDomain := DomainFromModel(userModel)
	return userDomain, nil
}
