package users_service

import (
	"context"
	"fmt"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	users_password "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/password"
)

type RegisterServiceResponse struct {
	User                  core_domain.User
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

func (s *usersService) RegisterUser(ctx context.Context, registerCommand core_domain.RegisterCommand) (RegisterServiceResponse, error) {
	if err := registerCommand.Validate(); err != nil {
		return RegisterServiceResponse{}, fmt.Errorf("validate user: %w", err)
	}

	hashedPassword, err := users_password.HashPassword(registerCommand.Password)
	if err != nil {
		return RegisterServiceResponse{}, fmt.Errorf("hash password: %w", err)
	}

	userDomain := core_domain.NewUserUninitialized(registerCommand.FullName, registerCommand.Email, hashedPassword)

	user, err := s.usersRepository.CreateUser(ctx, userDomain)
	if err != nil {
		return RegisterServiceResponse{}, fmt.Errorf("create user in repository: %w", err)
	}

	sessionServiceResponse, err := s.sessionsService.CreateSession(ctx, user.ID)
	if err != nil {
		return RegisterServiceResponse{
			User: user,
		}, fmt.Errorf("create session: %w", err)
	}

	return RegisterServiceResponse{
		User:                  user,
		AccessToken:           sessionServiceResponse.AccessToken,
		RefreshToken:          sessionServiceResponse.RefreshToken,
		AccessTokenExpiresAt:  sessionServiceResponse.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: sessionServiceResponse.RefreshTokenExpiresAt,
	}, nil
}
