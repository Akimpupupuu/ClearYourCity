package sessions_service

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
)

type SessionsService struct {
	sessionsRepository SessionsRepository
	tokenGenerator     *sessions_jwt.TokenGenerator
}

type SessionsRepository interface {
	CreateSession(ctx context.Context, session core_domain.Session) error
	RevokeSession(ctx context.Context, hashedRefreshToken string) error
	RevokeSessions(ctx context.Context, userID int) error
	GetSession(ctx context.Context, oldHashedToken string) (core_domain.Session, error)
}

func NewSessionsService(sessionsRepository SessionsRepository, tokenGenerator *sessions_jwt.TokenGenerator) *SessionsService {
	return &SessionsService{
		sessionsRepository: sessionsRepository,
		tokenGenerator:     tokenGenerator,
	}
}
