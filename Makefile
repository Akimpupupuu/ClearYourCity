include auth-service/.env
include task-service/.env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d auth-service-postgres task-service-postgres task-redis task-kafka

env-down:
	@docker compose down auth-service-postgres task-service-postgres task-redis task-kafka

env-cleanup:
	@docker compose down auth-service-postgres task-service-postgres task-redis task-kafka && \
	sudo rm -rf ${PROJECT_ROOT}/out/auth/pgdata/* && \
	sudo rm -rf ${PROJECT_ROOT}/out/task/pgdata/* && \
	sudo rm -rf ${PROJECT_ROOT}/out/task/kafkadata/* && \
	sudo rm -rf ${PROJECT_ROOT}/out/task/redisdata/*

migrate-create-auth:
	@if [ -z "$(seq)" ]; then \
		echo "Parametr seq is empty"; \
		exit 1; \
	fi;
	docker compose run --rm auth-migrate create -ext sql -dir /migrations -seq $(seq)

migrate-up-auth:
	@docker compose run --rm auth-migrate \
	-path /migrations \
	-database postgres://${AUTH_POSTGRES_USER}:${AUTH_POSTGRES_PASSWORD}@auth-service-postgres:5432/${AUTH_POSTGRES_DB}?sslmode=disable \
	up

migrate-down-auth:
	@docker compose run --rm auth-migrate \
	-path /migrations \
	-database postgres://${AUTH_POSTGRES_USER}:${AUTH_POSTGRES_PASSWORD}@auth-service-postgres:5432/${AUTH_POSTGRES_DB}?sslmode=disable \
	down

migrate-create-task:
	@if [ -z "$(seq)" ]; then \
		echo "Parametr seq is empty"; \
		exit 1; \
	fi;
	docker compose run --rm task-migrate create -ext sql -dir /migrations -seq $(seq)

migrate-up-task:
	@docker compose run --rm task-migrate \
	-path /migrations \
	-database postgres://${TASK_POSTGRES_USER}:${TASK_POSTGRES_PASSWORD}@task-service-postgres:5432/${TASK_POSTGRES_DB}?sslmode=disable \
	up

migrate-down-task:
	@docker compose run --rm task-migrate \
	-path /migrations \
	-database postgres://${TASK_POSTGRES_USER}:${TASK_POSTGRES_PASSWORD}@task-service-postgres:5432/${TASK_POSTGRES_DB}?sslmode=disable \
	down
	
swagger-gen-auth:
	@docker compose run --rm swagger-auth \
		init \
		-g cmd/auth/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

swagger-gen-task:
	@docker compose run --rm swagger-task \
		init \
		-g cmd/task/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

auth-run:
	@cd auth-service && \
	export AUTH_POSTGRES_HOST=localhost && \
	export AUTH_POSTGRES_PORT=5430 && \
	go run cmd/auth/main.go

auth-deploy:
	@docker compose up -d --build auth-service

auth-undeploy:
	@docker compose down auth-service

task-run:
	@cd task-service && \
	export TASK_POSTGRES_HOST=localhost && \
	export TASK_POSTGRES_PORT=5431 && \
	export TASK_KAFKA_BROKERS=localhost:9092 && \
	export TASK_REDIS_ADDR=localhost:6379 && \
	go run cmd/task/main.go

task-deploy:
	@docker compose up -d --build task-service

task-undeploy:
	@docker compose down task-service
