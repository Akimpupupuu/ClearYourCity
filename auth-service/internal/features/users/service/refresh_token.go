package users_service

import (
	"context"
	"fmt"
)

type RefreshTokenServiceResponse LoginServiceResponse

func (s *UsersService) RefreshToken(ctx context.Context, refreshToken string) (RefreshTokenServiceResponse, error) {
	sessionServiceResponse, err := s.sessionsService.RefreshToken(ctx, refreshToken)
	if err != nil {
		return RefreshTokenServiceResponse{}, fmt.Errorf("refresh token: %w", err)
	}

	return RefreshTokenServiceResponse{
		AccessToken:           sessionServiceResponse.AccessToken,
		RefreshToken:          sessionServiceResponse.RefreshToken,
		AccessTokenExpiresAt:  sessionServiceResponse.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: sessionServiceResponse.RefreshTokenExpiresAt,
	}, nil
}
