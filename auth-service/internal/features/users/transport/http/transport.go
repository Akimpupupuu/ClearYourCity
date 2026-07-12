package users_transport_http

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	users_service "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/service"
)

type UsersHandler struct {
	usersService UsersService
}

type UsersService interface {
	RegisterUser(ctx context.Context, cmd core_domain.RegisterCommand) (users_service.RegisterServiceResponse, error)
}

func NewUsersHandler(usersService UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}
