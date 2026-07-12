package sessions_service

import (
	"time"
)

type SessionServiceResponse struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}
