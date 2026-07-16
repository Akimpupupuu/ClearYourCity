package users_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	SELECT id, version, full_name, email, hashed_password, created_at
	FROM auth_service.users
	WHERE email=$1
	`

	row := r.pool.QueryRow(ctx, query, email)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.Email,
		&userModel.HashedPassword,
		&userModel.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf("user with email = '%s': %w", email, core_errors.ErrNotFound)
		}

		return core_domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	userDomain := DomainFromModel(userModel)
	return userDomain, nil
}
