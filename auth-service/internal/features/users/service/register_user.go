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
	Session               core_domain.Session
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

func (s *UsersService) RegisterUser(ctx context.Context, cmd core_domain.RegisterCommand) (RegisterServiceResponse, error) {
	if err := cmd.Validate(); err != nil {
		return RegisterServiceResponse{}, fmt.Errorf("validate user: %w", err)
	}

	hashedPassword, err := users_password.HashPassword(cmd.Password)
	if err != nil {
		return RegisterServiceResponse{}, fmt.Errorf("hash password: %w", err)
	}

	userDomain := core_domain.NewUserUninitialized(cmd.FullName, cmd.Email, hashedPassword)

	user, err := s.usersRepository.CreateUser(ctx, userDomain)
	if err != nil {
		return RegisterServiceResponse{}, fmt.Errorf("create user in repository: %w", err)
	}

	sessionServiceResponse, err := s.sessionsService.CreateSession(ctx, user.ID)
	if err != nil {
		return RegisterServiceResponse{
			User:                  user,
			Session:               core_domain.Session{},
			AccessToken:           "",
			RefreshToken:          "",
			AccessTokenExpiresAt:  time.Time{},
			RefreshTokenExpiresAt: time.Time{},
		}, fmt.Errorf("create session: %w", err)
	}

	return RegisterServiceResponse{
		User:                  user,
		Session:               sessionServiceResponse.Session,
		AccessToken:           sessionServiceResponse.AccessToken,
		RefreshToken:          sessionServiceResponse.RefreshToken,
		AccessTokenExpiresAt:  sessionServiceResponse.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: sessionServiceResponse.RefreshTokenExpiresAt,
	}, nil
}
