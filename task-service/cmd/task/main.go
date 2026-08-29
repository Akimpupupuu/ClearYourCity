package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	core_jwt "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/jwt"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	core_zap_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger/zap"
	core_postgres_pool "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/postgres/pool"
	core_redis "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/redis"
	http_router "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/router"
	http_server "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/server"
	core_kafka "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/kafka"
	tasks_outbox "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/outbox"
	tasks_kafka "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/repository/kafka"
	tasks_postgres "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/repository/postgres"
	tasks_redis "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/repository/redis"
	tasks_service "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/service"
	tasks_transport_http "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/transport/http"
	"github.com/go-chi/chi"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("critical application error: %v", err)
	}
}

func run() error {
	const (
		apiVersionV1 = "v1"
	)

	time.Local = time.UTC

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger, err := core_zap_logger.NewLogger()
	if err != nil {
		return fmt.Errorf("failed init application logger: %w", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		return fmt.Errorf("failed to init postgres connection pool: %w", err)
	}
	defer pool.Close()

	// logger.Debug("initializing transaction manager")
	// transactionManager := core_postgres_transaction.NewTransactionManager(pool.Pool)

	logger.Debug("initializing token generator")
	tokenGenerator := core_jwt.NewTokenGenerator(core_jwt.NewConfigMust())

	logger.Debug("initializing redis database")
	redis, err := core_redis.NewRedisClient(ctx, core_redis.NewConfigMust())
	if err != nil {
		return fmt.Errorf("failed to init redis: %w", err)
	}
	defer func() {
		_ = redis.Close()
	}()

	logger.Debug("initializing redis repository")
	redisRepository := tasks_redis.NewTasksRedis(redis)

	logger.Debug("initializing kafka producer")
	producer := core_kafka.NewProducer(core_kafka.NewConfigMust())
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error("failed to close kafka producer", core_logger.Err(err))
		}
	}()

	logger.Debug("initializing producer repository")
	producerRepository := tasks_kafka.NewTasksKafkaRepository(producer)

	logger.Debug("initializing tasks feature")
	tasksRepository := tasks_postgres.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository, redisRepository, logger)
	tasksTransportHTTP := tasks_transport_http.NewTasksHandler(tasksService, tokenGenerator)

	logger.Debug("initializing router")
	router := http_router.NewRouter(logger)
	router.RouteApi(apiVersionV1, func(apiRouter chi.Router) {
		tasksTransportHTTP.Register(apiRouter)
	})

	logger.Debug("initializing outbox worker")
	worker := tasks_outbox.NewWorker(tasksRepository, redisRepository, producerRepository, logger, tasks_outbox.NewConfigMust())
	go worker.Start(ctx)

	logger.Debug("initializing HTTP Server")
	server := http_server.NewHTTPServer(router, http_server.NewConfigMust(), logger)
	if err = server.Run(ctx); err != nil {
		return fmt.Errorf("run http server: %w", err)
	}

	return nil
}
