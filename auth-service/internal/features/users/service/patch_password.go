package users_service

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
	users_password "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/password"
)

func (s *usersService) PatchPassword(ctx context.Context, oldPassword string, newPassword string) error {
	if err := validatePassword(newPassword); err != nil {
		return fmt.Errorf("validate password: %w", err)
	}

	claims, ok := sessions_jwt.FromContext(ctx)
	if !ok {
		return fmt.Errorf("get claims from context: %w", core_errors.ErrUnauthorized)
	}

	user, err := s.usersRepository.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return fmt.Errorf("get user from repository: %v: %w", err, core_errors.ErrUnauthorized)
		}

		return fmt.Errorf("get user from repository: %w", err)
	}

	if err := users_password.VerifyPassword(user.HashedPassword, oldPassword); err != nil {
		return fmt.Errorf("user authentication failed: %w", core_errors.ErrUnauthorized)
	}

	hashedNewPassword, err := users_password.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user.HashedPassword = hashedNewPassword

	if err := s.usersRepository.PatchPassword(ctx, user); err != nil {
		return fmt.Errorf("patch password in repository: %w", err)
	}

	if err := s.sessionsService.RevokeSessions(ctx, claims.UserID); err != nil {
		return fmt.Errorf("revoke sessions: %w", err)
	}

	return nil
}

func validatePassword(password string) error {
	passwordLength := len([]rune(password))
	if passwordLength < 8 {
		return fmt.Errorf("invalid password length: %d: %w", passwordLength, core_errors.ErrInvalidArgument)
	}

	return nil
}
