package users_transport_http

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
)

type UsersHandler struct {
	usersService UsersService
}

type UsersService interface {
	RegisterUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
}

func NewUsersHandler(usersService UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}
