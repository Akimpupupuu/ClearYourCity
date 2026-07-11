package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	http_router "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/http/router"
	http_server "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/http/server"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	core_postgres_pool "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/postgres/pool"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
	"go.uber.org/zap"
)

func main() {
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

	// repo level
	// service level
	// transport level

	logger.Debug("initializing router")
	router := http_router.NewRouter(logger, tokenGenerator)
	router.RegisterRoutes()

	logger.Debug("initializing HTTP Server")
	server := http_server.NewHTTPServer(router, http_server.NewConfigMust(), logger)
	server.Run(ctx)
}
