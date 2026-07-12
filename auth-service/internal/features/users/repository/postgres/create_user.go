package users_postrges_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgxViolatesUniqueErrorCode = "23505"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	INSERT INTO auth_service.users (full_name, email, hashed_password, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, version, full_name, email, hashed_password, created_at;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.Email, user.HashedPassword, user.CreatedAt)

	var usersModel UserModel
	if err := row.Scan(
		&usersModel.ID,
		&usersModel.Version,
		&usersModel.FullName,
		&usersModel.Email,
		&usersModel.HashedPassword,
		&usersModel.CreatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgxViolatesUniqueErrorCode {
				return core_domain.User{}, fmt.Errorf("%v: user with email = '%s': %w", err, user.Email, core_errors.ErrConflict)
			}
		}

		return core_domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := core_domain.NewUser(
		usersModel.ID,
		usersModel.Version,
		usersModel.FullName,
		usersModel.Email,
		usersModel.HashedPassword,
		usersModel.CreatedAt,
	)

	return userDomain, nil
}
