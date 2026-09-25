package http_router

import (
	"net/http"
	"path"

	"github.com/Akimpupupuu/ClearYourCity/auth-service/docs"
	http_middleware "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/middleware"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter(log *zap.Logger) *Router {
	r := chi.NewRouter()

	r.Use(http_middleware.CORS())
	r.Use(http_middleware.RequestID())
	r.Use(http_middleware.Logger(log))
	r.Use(http_middleware.Trace())
	r.Use(http_middleware.Panic())

	router := &Router{mux: r}
	router.RegisterSwagger()

	return router
}

func (r *Router) RouteApi(version string, builder func(apiRouter chi.Router)) {
	path := path.Join("/api/", version)
	r.mux.Route(path, func(router chi.Router) {
		builder(router)
	})
}

func (r *Router) RegisterSwagger() {
	r.mux.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"), httpSwagger.DefaultModelsExpandDepth(-1)))
	r.mux.Get("/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (r *Router) GetMux() http.Handler {
	return r.mux
}
