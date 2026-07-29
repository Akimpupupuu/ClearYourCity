package http_router

import (
	"net/http"
	"path"

	http_middleware "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/middleware"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter(log *zap.Logger) *Router {
	r := chi.NewRouter()

	r.Use(http_middleware.RequestID())
	r.Use(http_middleware.Logger(log))
	r.Use(http_middleware.Trace())
	r.Use(http_middleware.Panic())

	return &Router{mux: r}
}

func (r *Router) RouteApi(version string, builder func(apiRouter chi.Router)) {
	path := path.Join("/api/", version)
	r.mux.Route(path, func(router chi.Router) {
		builder(router)
	})
}

func (r *Router) GetMux() http.Handler {
	return r.mux
}
