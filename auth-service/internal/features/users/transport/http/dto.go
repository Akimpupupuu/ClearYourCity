package users_transport_http

import (
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
)

type ResponseUserDTO struct {
	ID        int       `json:"id" example:"1"`
	Version   int       `json:"version" example:"2"`
	FullName  string    `json:"full_name" example:"Иван"`
	Email     string    `json:"email" example:"ivan@gmail.com"`
	CreatedAt time.Time `json:"created_at" example:"2026-08-29T18:51:08.085831Z"`
}

type ResponseRegisterDTO struct {
	AccessToken          string          `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkFraW0iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTgyNjIzOTAyMn0.bFp3dThEamI3SjFza0FFbF9MUDJiWDF5YmNLdzVVY2E1eXQ4N0xfNDJnNA"`
	AccessTokenExpiresAt time.Time       `json:"access_token_expires_at" example:"2026-08-29T18:51:08.085831Z"`
	User                 ResponseUserDTO `json:"user"`
}

type ResponseLoginDTO struct {
	AccessToken          string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkFraW0iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTgyNjIzOTAyMn0.bFp3dThEamI3SjFza0FFbF9MUDJiWDF5YmNLdzVVY2E1eXQ4N0xfNDJnNA"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at" example:"2026-08-29T18:51:08.085831Z"`
}

func UserDTOFromDomain(user *core_domain.User) ResponseUserDTO {
	return ResponseUserDTO{
		ID:        user.ID,
		Version:   user.Version,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func RegisterDTOFromService(user *core_domain.User, accessToken string, accessTokenExpiresAt time.Time) ResponseRegisterDTO {
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
