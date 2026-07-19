package users_service

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
)

func (s *usersService) GetUser(ctx context.Context, userID int) (core_domain.User, error) {
	user, err := s.usersRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return core_domain.User{}, fmt.Errorf("get user from repository: %v: %w", err, core_errors.ErrUnauthorized)
		}

		return core_domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
