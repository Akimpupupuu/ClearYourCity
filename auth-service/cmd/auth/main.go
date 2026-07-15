package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	core_postgres_pool "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/postgres/pool"
	http_router "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/router"
	http_server "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/server"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
	sessions_repository "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/repository/postgres"
	sessions_service "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/service"
	users_repository "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/repository/postgres"
	users_service "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/service"
	users_transport_http "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/users/transport/http"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	time.Local = time.UTC

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger, err := core_logger.InitLogger()
	if err != nil {
		fmt.Println("failed init application logger:", err)
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	logger.Debug("initializing token generator")
	tokenGenerator := sessions_jwt.NewTokenGenerator(sessions_jwt.NewConfigMust())

	logger.Debug("initializing sessions feature")
	sessionsRepository := sessions_repository.NewSessionsRepository(pool)
	sessionsService := sessions_service.NewSessionsService(sessionsRepository, tokenGenerator)

	logger.Debug("initializing sessions feature")
	usersRepository := users_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository, sessionsService)
	usersTransportHTTP := users_transport_http.NewUsersHandler(usersService, tokenGenerator)

	logger.Debug("initializing router")
	router := http_router.NewRouter(logger)
	router.RouteApi("v1", func(apiRouter chi.Router) {
		usersTransportHTTP.Register(apiRouter)
	})

	logger.Debug("initializing HTTP Server")
	server := http_server.NewHTTPServer(router, http_server.NewConfigMust(), logger)
	server.Run(ctx)
}
