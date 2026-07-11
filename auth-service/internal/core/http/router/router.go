package http_router

import (
	"net/http"

	http_middleware "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/http/middleware"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Router struct {
	mux            *chi.Mux
	log            *zap.Logger
	tokenGenerator *sessions_jwt.TokenGenerator
}

func NewRouter(log *zap.Logger, tokenGenerator *sessions_jwt.TokenGenerator) *Router {
	return &Router{
		mux:            chi.NewRouter(),
		log:            log,
		tokenGenerator: tokenGenerator,
	}
}

func (r *Router) RegisterRoutes() {

	r.mux.Use(http_middleware.RequestID())
	r.mux.Use(http_middleware.Logger(r.log))
	r.mux.Use(http_middleware.Trace())
	r.mux.Use(http_middleware.Panic())

	r.mux.Route("/api/v1", func(r chi.Router) {

	})
}

func (r *Router) GetMux() http.Handler {
	return r.mux
}
