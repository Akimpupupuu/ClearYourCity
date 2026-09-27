package tasks_redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const (
	acceptTTL = time.Hour * 168
)

func (r *tasksRedis) CreateAction(ctx context.Context, token string, taskID int) error {
	redisKey := fmt.Sprintf(tokenPrefix, token)

	if err := r.redis.Set(ctx, redisKey, strconv.Itoa(taskID), acceptTTL).Err(); err != nil {
		return fmt.Errorf("set value in redis: %w", err)
	}

	return nil
}
