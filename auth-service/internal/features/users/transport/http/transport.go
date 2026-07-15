package users_transport_http

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
	users_service "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/service"
	"github.com/go-chi/chi"
)

type UsersHandler struct {
	usersService   UsersService
	tokenGenerator *sessions_jwt.TokenGenerator
}

type UsersService interface {
	RegisterUser(ctx context.Context, cmd core_domain.RegisterCommand) (users_service.RegisterServiceResponse, error)
	LoginUser(ctx context.Context, loginCommand core_domain.LoginCommand) (users_service.LoginServiceResponse, error)
	LogoutUser(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (users_service.RefreshTokenServiceResponse, error)
}

func NewUsersHandler(usersService UsersService, tokenGenerator *sessions_jwt.TokenGenerator) *UsersHandler {
	return &UsersHandler{
		usersService:   usersService,
		tokenGenerator: tokenGenerator,
	}
}

func (h *UsersHandler) Register(router chi.Router) {
	router.Route("/auth", func(subRouter chi.Router) {
		subRouter.Post("/register", h.RegisterUser)
		subRouter.Post("/login", h.LoginUser)
		subRouter.Post("/refresh", h.RefreshToken)
		subRouter.Post("/logout", h.LogoutUser)
	})
}
