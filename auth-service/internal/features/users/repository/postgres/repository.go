package users_repository

import (
	core_postgres_pool "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/postgres/pool"
)

type usersRepository struct {
	pool *core_postgres_pool.Pool
}

func NewUsersRepository(pool *core_postgres_pool.Pool) *usersRepository {
	return &usersRepository{
		pool: pool,
	}
}
