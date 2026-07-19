package users_service

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
)

func (s *usersService) PatchUser(ctx context.Context, fullName *string, email *string) (core_domain.User, error) {
	claims, ok := sessions_jwt.FromContext(ctx)
	if !ok {
		return core_domain.User{}, fmt.Errorf("get claims from context: %w", core_errors.ErrUnauthorized)
	}

	user, err := s.usersRepository.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return core_domain.User{}, fmt.Errorf("get user from repository: %v: %w", err, core_errors.ErrUnauthorized)
		}

		return core_domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	if err := user.ApplyPatch(fullName, email); err != nil {
		return core_domain.User{}, fmt.Errorf("apply patch: %w", err)
	}

	user, err = s.usersRepository.PatchUser(ctx, user)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("patch user in repository: %w", err)
	}

	return user, nil

}
