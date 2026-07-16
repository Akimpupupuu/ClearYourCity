package users_service

import (
	"context"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
)

func (s *UsersService) GetUser(ctx context.Context, userID int) (core_domain.User, error) {
	user, err := s.usersRepository.GetUserByID(ctx, userID)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
