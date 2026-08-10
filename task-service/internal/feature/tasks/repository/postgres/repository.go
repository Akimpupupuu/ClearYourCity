package tasks_postgres

import core_postgres_pool "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/postgres/pool"

type tasksRepository struct {
	pool *core_postgres_pool.Pool
}

func NewTasksRepository(pool *core_postgres_pool.Pool) *tasksRepository {
	return &tasksRepository{
		pool: pool,
	}
}
