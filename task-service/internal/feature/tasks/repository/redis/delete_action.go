package tasks_redis

import (
	"context"
	"fmt"
)

func (r *tasksRedis) DeleteRedis(ctx context.Context, token string) error {
	redisKey := fmt.Sprintf(tokenPrefix, token)

	if err := r.redis.Del(ctx, redisKey).Err(); err != nil {
		return fmt.Errorf("delete action from redis: %w", err)
	}

	return nil
}
