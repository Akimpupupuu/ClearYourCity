package users_transport_http

import (
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
)

type ResponseUserDTO struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseRegisterDTO struct {
	AccessToken          string          `json:"access_token"`
	AccessTokenExpiresAt time.Time       `json:"access_token_expires_at"`
	User                 ResponseUserDTO `json:"user"`
}

type ResponseLoginDTO struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func UserDTOFromDomain(user core_domain.User) ResponseUserDTO {
	return ResponseUserDTO{
		ID:        user.ID,
		Version:   user.Version,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func RegisterDTOFromService(user core_domain.User, accessToken string, accessTokenExpiresAt time.Time) ResponseRegisterDTO {
	return ResponseRegisterDTO{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
		User: ResponseUserDTO{
			ID:        user.ID,
			Version:   user.Version,
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}
}

func LoginDTOFromService(accessToken string, accessTokenExpiresAt time.Time) ResponseLoginDTO {
	return ResponseLoginDTO{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	}
}
