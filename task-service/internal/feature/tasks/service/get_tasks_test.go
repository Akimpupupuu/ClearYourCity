package tasks_service

import (
	"context"
	"fmt"
	"testing"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_zap_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger/zap"
	mock_tasks_service "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testGetTasksInput struct {
	name        string
	userID      int
	limit       int
	offset      int
	dbExec      bool
	expResultDb []*core_domain.Task
	dbErr       error
	expErr      error
}

func TestGetTasks(t *testing.T) {
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)

	const (
		taskIDFirst   = 1
		taskIDSecond  = 2
		userID        = 1
		version       = 1
		limit         = 20
		invalidLimit  = 101
		offset        = 0
		invalidOffset = -1
	)

	outputDbTasks := []*core_domain.Task{
		core_domain.NewTask(taskIDFirst, version, userID, "Clean the park", "Remove plastic bottles near the park", core_domain.StatusCreated, time.Date(2026, time.September, 28, 18, 0, 0, 0, time.UTC), nil),
		core_domain.NewTask(taskIDSecond, version, userID, "Clean the lake", "Remove plastic bottles near the lake", core_domain.StatusCreated, time.Date(2026, time.September, 28, 15, 0, 0, 0, time.UTC), nil),
	}

	inputs := []testGetTasksInput{
		{
			"success: got",
			userID,
			limit,
			offset,
			true,
			outputDbTasks,
			nil,
			nil,
		},
		{
			"fail: validation error: limit > 100",
			userID,
			invalidLimit,
			offset,
			false,
			nil,
			nil,
			fmt.Errorf("validate pagination values: 'limit' must be non-negative value and smaller then 100: 101: invalid argument"),
		},
		{
			"fail: validation error: offset < 0",
			userID,
			limit,
			invalidOffset,
			false,
			nil,
			nil,
			fmt.Errorf("validate pagination values: 'offset' must be non-negative value: -1: invalid argument"),
		},
		{
			"fail: database error",
			userID,
			limit,
			offset,
			true,
			nil,
			fmt.Errorf("sql: connection refused"),
			fmt.Errorf("get tasks from repository: sql: connection refused"),
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			ctl := gomock.NewController(t)

			ctx := context.Background()
			postgres := mock_tasks_service.NewMockTasksRepository(ctl)
			redis := mock_tasks_service.NewMockTasksRedis(ctl)

			if input.dbExec {
				postgres.EXPECT().GetTasks(ctx, input.userID, input.limit, input.offset).Return(input.expResultDb, input.dbErr)
			}

			service := NewTasksService(postgres, redis, log)
			tasks, err := service.GetTasks(ctx, input.userID, input.limit, input.offset)
			if input.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, input.expErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, input.expResultDb, tasks)
		})
	}
}
