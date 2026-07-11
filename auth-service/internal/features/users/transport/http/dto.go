package users_transport_http

import (
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
)

type ResponseUserDTO struct {
	ID        int
	Version   int
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func DTOFromDomain(domain core_domain.User) ResponseUserDTO {
	return ResponseUserDTO{
		ID:        domain.ID,
		Version:   domain.Version,
		FullName:  domain.FullName,
		Email:     domain.Email,
		Password:  domain.Password,
		CreatedAt: domain.CreatedAt,
	}
}
