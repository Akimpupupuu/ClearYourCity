package tasks_redis

import (
	"context"
	"fmt"
	"time"
)

const (
	PerformTTL = time.Hour * 24 * 30
)

func (r *tasksRedis) ProlongAction(ctx context.Context, token string) error {
	redisKey := fmt.Sprintf(tokenPrefix, token)

	if err := r.redis.Expire(ctx, redisKey, PerformTTL).Err(); err != nil {
		return fmt.Errorf("prolong token in redis: %w", err)
	}

	return nil
}
