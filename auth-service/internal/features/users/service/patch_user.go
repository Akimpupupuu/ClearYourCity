package users_service

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
)

func (s *usersService) PatchUser(ctx context.Context, userID int, patchUserCommand core_domain.PatchUserCommand) (*core_domain.User, error) {
	if err := patchUserCommand.Validate(); err != nil {
		return nil, fmt.Errorf("validate 'patchUserCommand': %w", err)
	}

	user, err := s.usersRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return nil, fmt.Errorf("get user from repository: %v: %w", err, core_errors.ErrUnauthorized)
		}

		return nil, fmt.Errorf("get user from repository: %w", err)
	}

	if err := user.ApplyPatch(patchUserCommand.FullName, patchUserCommand.Email); err != nil {
		return nil, fmt.Errorf("apply patch: %w", err)
	}

	user, err = s.usersRepository.PatchUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("patch user in repository: %w", err)
	}

	return user, nil

}
