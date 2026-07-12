package users_repository

import (
	core_postgres_pool "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/postgres/pool"
)

type UsersRepository struct {
	pool *core_postgres_pool.Pool
}

func NewUsersRepository(pool *core_postgres_pool.Pool) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
