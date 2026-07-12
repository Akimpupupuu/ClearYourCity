package users_transport_http

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	users_service "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/service"
	"github.com/go-chi/chi"
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

func (h *UsersHandler) Register(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
	})
}
