package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	core_zap_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger/zap"
	core_postgres_pool "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/postgres/pool"
	core_redis "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/redis"
	http_router "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/router"
	http_server "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/server"
	core_kafka "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/kafka"
	tasks_kafka "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/repository/kafka"
	tasks_postgres "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/repository/postgres"
	tasks_redis "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/repository/redis"
	tasks_service "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/service"
	tasks_transport_http "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/transport/http"
	"github.com/go-chi/chi"
)

func main() {
	const (
		apiVersionV1 = "v1"
	)

	time.Local = time.UTC

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger, err := core_zap_logger.NewLogger()
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
		logger.Fatal("failed to init postgres connection pool", core_logger.Err(err))
	}
	defer pool.Close()

	// logger.Debug("initializing transaction manager")
	// transactionManager := core_postgres_transaction.NewTransactionManager(pool.Pool)

	// logger.Debug("initializing token generator")
	// tokenGenerator := core_jwt.NewTokenGenerator(core_jwt.NewConfigMust())

	logger.Debug("initializing redis database")
	redis, err := core_redis.NewRedisClient(ctx, core_redis.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init redis", core_logger.Err(err))
	}

	logger.Debug("initializing redis repository")
	redisRepository := tasks_redis.NewTasksRedis(redis)

	logger.Debug("initializing kafka producer")
	producer := core_kafka.NewProducer(core_kafka.NewConfigMust())

	logger.Debug("initializing producer repository")
	producerRepository := tasks_kafka.NewTasksKafkaRepository(producer)

	logger.Debug("initializing tasks feature")
	tasksRepository := tasks_postgres.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository, redisRepository, producerRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHandler(tasksService)

	logger.Debug("initializing router")
	router := http_router.NewRouter(logger)
	router.RouteApi(apiVersionV1, func(apiRouter chi.Router) {
		tasksTransportHTTP.Register(apiRouter)
	})

	logger.Debug("initializing HTTP Server")
	server := http_server.NewHTTPServer(router, http_server.NewConfigMust(), logger)
	server.Run(ctx)

	logger.Debug("waiting for background tasks to complete...")
	tasksService.GracefulShutdown()

	logger.Debug("application entirely stopped")
}
