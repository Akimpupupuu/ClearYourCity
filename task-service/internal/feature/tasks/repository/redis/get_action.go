package tasks_redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	"github.com/redis/go-redis/v9"
)

func (r *tasksRedis) GetAction(ctx context.Context, token string) (int, error) {
	redisKey := fmt.Sprintf(tokenPrefix, token)

	val, err := r.redis.Get(ctx, redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, fmt.Errorf("get value from redis: %v: %w", err, core_errors.ErrNotFound)
		}
		return 0, fmt.Errorf("get value from redis: %w", err)
	}

	taskID, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("parse taskID from redis: %w", err)
	}

	return taskID, nil
}
