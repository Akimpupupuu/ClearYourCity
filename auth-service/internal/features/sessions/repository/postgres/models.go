package sessions_repository

import (
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
)

type SessionModel struct {
	ID               string
	UserID           int
	RefreshTokenHash string
	IsRevoked        bool
	CreatedAt        time.Time
	ExpiresAt        time.Time
}

func DomainFromModel(model SessionModel) core_domain.Session {
	return core_domain.NewSession(
		model.ID,
		model.UserID,
		model.RefreshTokenHash,
		model.IsRevoked,
		model.CreatedAt,
		model.ExpiresAt,
	)
}
