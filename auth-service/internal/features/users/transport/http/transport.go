package users_transport_http

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	http_middleware "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/middleware"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
	users_service "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/service"
	"github.com/go-chi/chi"
)

type usersHandler struct {
	usersService   UsersService
	tokenGenerator *sessions_jwt.TokenGenerator
}

type UsersService interface {
	RegisterUser(ctx context.Context, cmd core_domain.RegisterCommand) (users_service.RegisterServiceResponse, error)
	LoginUser(ctx context.Context, loginCommand core_domain.LoginCommand) (users_service.LoginServiceResponse, error)
	LogoutUser(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (users_service.RefreshTokenServiceResponse, error)
	GetUser(ctx context.Context, userID int) (core_domain.User, error)
	PatchUser(ctx context.Context, fullName *string, email *string) (core_domain.User, error)
	PatchPassword(ctx context.Context, oldPassword string, newPassword string) error
}

func NewUsersHandler(usersService UsersService, tokenGenerator *sessions_jwt.TokenGenerator) *usersHandler {
	return &usersHandler{
		usersService:   usersService,
		tokenGenerator: tokenGenerator,
	}
}

func (h *usersHandler) Register(router chi.Router) {
	router.Route("/auth", func(subRouter chi.Router) {
		subRouter.Post("/register", h.RegisterUser)
		subRouter.Post("/login", h.LoginUser)
		subRouter.Post("/refresh", h.RefreshToken)
		subRouter.Post("/logout", h.LogoutUser)

		subRouter.Group(func(protected_router chi.Router) {
			protected_router.Use(http_middleware.Auth(h.tokenGenerator))

			protected_router.Get("/get", h.GetUser)
			protected_router.Patch("/patch_password", h.PatchPassword)
			protected_router.Patch("/patch_user", h.PatchUser)
		})
	})
}
