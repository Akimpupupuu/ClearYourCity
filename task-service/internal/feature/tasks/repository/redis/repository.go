package tasks_redis

import (
	core_redis "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/redis"
)

var (
	tokenPrefix = "task_action:%s"
)

type tasksRedis struct {
	redis *core_redis.Redis
}

func NewTasksRedis(redis *core_redis.Redis) *tasksRedis {
	return &tasksRedis{
		redis: redis,
	}
}
