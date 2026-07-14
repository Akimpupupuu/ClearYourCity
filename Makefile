include auth-service/.env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d auth-service-postgres

env-down:
	@docker compose down auth-service-postgres

env-cleanup:
	@docker compose down auth-service-postgres && \
	sudo rm -rf ${PROJECT_ROOT}/out/auth/pgdata

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Parametr seq is empty"; \
		exit 1; \
	fi;
	docker compose run --rm auth-migrate create -ext sql -dir /migrations -seq $(seq)

migrate-up:
	@docker compose run --rm auth-migrate \
	-path /migrations \
	-database postgres://${AUTH_POSTGRES_USER}:${AUTH_POSTGRES_PASSWORD}@auth-service-postgres:5432/${AUTH_POSTGRES_DB}?sslmode=disable \
	up

migrate-down:
	@docker compose run --rm auth-migrate \
	-path /migrations \
	-database postgres://${AUTH_POSTGRES_USER}:${AUTH_POSTGRES_PASSWORD}@auth-service-postgres:5432/${AUTH_POSTGRES_DB}?sslmode=disable \
	down

auth-run:
	@cd auth-service && go run cmd/auth/main.go
