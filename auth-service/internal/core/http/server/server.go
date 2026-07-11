package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	http_router "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/http/router"
	"go.uber.org/zap"
)

type HTTPServer struct {
	router *http_router.Router
	config Config
	log    *zap.Logger
}

func NewHTTPServer(router *http_router.Router, config Config, log *zap.Logger) *HTTPServer {
	return &HTTPServer{
		router: router,
		config: config,
		log:    log,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := http.Server{
		Addr:    s.config.Addr,
		Handler: s.router.GetMux(),
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownContext, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		err := server.Shutdown(shutdownContext)
		if err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
