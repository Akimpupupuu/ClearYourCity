package users_service

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	sessions_service "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/service"
)

type UsersService struct {
	usersRepository UsersRepository
	sessionsService SessionsService
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
	GetUser(ctx context.Context, email string) (core_domain.User, error)
}

type SessionsService interface {
	CreateSession(ctx context.Context, userID int) (sessions_service.SessionServiceResponse, error)
	RefreshToken(ctx context.Context, oldRefreshToken string) (sessions_service.SessionServiceResponse, error)
	RevokeSession(ctx context.Context, refreshToken string) error
	RevokeSessions(ctx context.Context, userID int) error
}

func NewUsersService(usersRepository UsersRepository, sessionsService SessionsService) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		sessionsService: sessionsService,
	}
}
