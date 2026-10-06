package tasks_service

import (
	"context"
	"fmt"
	"testing"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_zap_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger/zap"
	mock_tasks_service "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// it is better to use setup mocks
type testCreateTaskInput struct {
	name        string
	input       *core_domain.Task
	dbExec      bool
	expResultDb *core_domain.Task
	dbErr       error
	expErr      error
}

func TestCreateTask(t *testing.T) {
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)

	const (
		taskID  = 1
		userID  = 1
		version = 1
	)

	inputTask := core_domain.NewTaskUninitialized("Clean the park", "Remove plastic bottles near the lake", userID)
	outputDbTask := core_domain.NewTask(taskID, version, userID, "Clean the park", "Remove plastic bottles near the lake", core_domain.StatusCreated, inputTask.CreatedAt, nil)
	invalidTask := core_domain.NewTaskUninitialized("Clean the park", "plastic", userID)

	inputs := []testCreateTaskInput{
		{
			"success: created",
			inputTask,
			true,
			outputDbTask,
			nil,
			nil,
		},
		{
			"fail: validation error",
			invalidTask,
			false,
			nil,
			nil,
			fmt.Errorf("validate task: invalid 'description' length: 7: invalid argument"),
		},
		{
			"fail: database error",
			inputTask,
			true,
			nil,
			fmt.Errorf("sql: connection refused"),
			fmt.Errorf("create task in repository: sql: connection refused"),
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			ctl := gomock.NewController(t)

			ctx := context.Background()
			postgres := mock_tasks_service.NewMockTasksRepository(ctl)
			redis := mock_tasks_service.NewMockTasksRedis(ctl)

			if input.dbExec {
				postgres.EXPECT().CreateTask(ctx, input.input, gomock.Any()).Return(input.expResultDb, input.dbErr)
			}

			service := NewTasksService(postgres, redis, log)
			task, err := service.CreateTask(ctx, input.input)
			if input.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, input.expErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, input.expResultDb, task)
		})
	}
}
