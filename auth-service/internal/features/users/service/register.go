package users_service

import (
	"context"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	users_password "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/password"
)

func (s *UsersService) RegisterUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	// make validation
	if err := user.Validate(); err != nil {
		return core_domain.User{}, fmt.Errorf("validate user: %w", err)
	}

	hashedPassword, err := users_password.HashPassword(user.Password)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	user.Password = hashedPassword

	// create user
	// create session
	// call repo
	// return result

	return core_domain.User{}, nil
}
