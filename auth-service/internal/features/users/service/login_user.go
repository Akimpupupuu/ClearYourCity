package users_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	users_password "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/password"
)

type LoginServiceResponse struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

func (s *UsersService) LoginUser(ctx context.Context, loginCommand core_domain.LoginCommand) (LoginServiceResponse, error) {
	user, err := s.usersRepository.GetUser(ctx, loginCommand.Email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return LoginServiceResponse{}, fmt.Errorf("user authentication failed: %w", core_errors.ErrUnauthorized)
		}

		return LoginServiceResponse{}, fmt.Errorf("get user from repository: %w", err)
	}

	if err := users_password.VerifyPassword(user.HashedPassword, loginCommand.Password); err != nil {
		return LoginServiceResponse{}, fmt.Errorf("user authentication failed: %w", core_errors.ErrUnauthorized)
	}

	sessionServiceResponse, err := s.sessionsService.CreateSession(ctx, user.ID)
	if err != nil {
		return LoginServiceResponse{}, fmt.Errorf("create session: %w", err)
	}

	return LoginServiceResponse{
		AccessToken:           sessionServiceResponse.AccessToken,
		RefreshToken:          sessionServiceResponse.RefreshToken,
		AccessTokenExpiresAt:  sessionServiceResponse.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: sessionServiceResponse.RefreshTokenExpiresAt,
	}, nil
}
