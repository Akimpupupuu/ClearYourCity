package sessions_repository

import core_postgres_pool "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/postgres/pool"

type sessionsRepository struct {
	pool *core_postgres_pool.Pool
}

func NewSessionsRepository(pool *core_postgres_pool.Pool) *sessionsRepository {
	return &sessionsRepository{
		pool: pool,
	}
}
