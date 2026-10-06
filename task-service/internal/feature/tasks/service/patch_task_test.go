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

type testPatchTaskInput struct {
	name      string
	taskID    int
	userID    int
	command   core_domain.PatchTaskCommand
	mocks     func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task)
	expResult *core_domain.Task
	expErr    error
}

func TestPatchTask(t *testing.T) {
	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)

	const (
		taskID        = 1
		userID        = 1
		invalidUserID = 2
		version       = 1
	)

	title := "Clean the lake"
	invalidTitle := "Hi"
	description := "Remove plastic bottles at the park, near the lake"
	invalidDescription := "plastic"

	initialTask := core_domain.NewTask(taskID, version, userID, "Clean the park", "Remove plastic bottles near the park", core_domain.StatusCreated, time.Date(2026, time.September, 28, 18, 0, 0, 0, time.UTC), nil)
	updatedTaskTitle := core_domain.NewTask(taskID, version, userID, "Clean the lake", "Remove plastic bottles near the park", core_domain.StatusCreated, time.Date(2026, time.September, 28, 18, 0, 0, 0, time.UTC), nil)
	updatedTaskDescription := core_domain.NewTask(taskID, version, userID, "Clean the park", "Remove plastic bottles at the park, near the lake", core_domain.StatusCreated, time.Date(2026, time.September, 28, 18, 0, 0, 0, time.UTC), nil)
	updatedTaskTitleAndDescription := core_domain.NewTask(taskID, version, userID, "Clean the lake", "Remove plastic bottles near the lake and park", core_domain.StatusCreated, time.Date(2026, time.September, 28, 18, 0, 0, 0, time.UTC), nil)

	inputs := []testPatchTaskInput{
		{
			"success: patched title",
			taskID,
			userID,
			core_domain.NewPatchTaskCommand(&title, nil),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchTask(ctx, gomock.Any()).Return(updatedTaskTitle, nil)
			},
			updatedTaskTitle,
			nil,
		},
		{
			"success: patched description",
			taskID,
			userID,
			core_domain.NewPatchTaskCommand(nil, &description),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchTask(ctx, gomock.Any()).Return(updatedTaskDescription, nil)
			},
			updatedTaskDescription,
			nil,
		},
		{
			"success: patched title and description",
			taskID,
			userID,
			core_domain.NewPatchTaskCommand(&title, &description),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchTask(ctx, gomock.Any()).Return(updatedTaskTitleAndDescription, nil)
			},
			updatedTaskTitleAndDescription,
			nil,
		},
		{
			"fail: get task from repository error",
			taskID,
			userID,
			core_domain.NewPatchTaskCommand(&title, &description),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				repo.EXPECT().GetTask(ctx, taskID).Return(nil, fmt.Errorf("sql: connection refused"))
			},
			nil,
			fmt.Errorf("get task from repository: sql: connection refused"),
		},
		{
			"fail: invalid user id",
			taskID,
			invalidUserID,
			core_domain.NewPatchTaskCommand(&title, &description),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
			},
			nil,
			fmt.Errorf("access denied: not found"),
		},
		{
			"fail: apply patch error",
			taskID,
			userID,
			core_domain.NewPatchTaskCommand(&invalidTitle, &invalidDescription),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
			},
			nil,
			fmt.Errorf("apply patch: validate task: invalid 'title' length: 2: invalid argument"),
		},
		{
			"fail: patch task in repository error",
			taskID,
			userID,
			core_domain.NewPatchTaskCommand(&title, &description),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchTask(ctx, gomock.Any()).Return(nil, fmt.Errorf("sql: connection refused"))
			},
			nil,
			fmt.Errorf("patch task in repository: sql: connection refused"),
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			ctl := gomock.NewController(t)

			postgres := mock_tasks_service.NewMockTasksRepository(ctl)
			redis := mock_tasks_service.NewMockTasksRedis(ctl)

			initialTaskCopy := *initialTask

			input.mocks(postgres, redis, &initialTaskCopy)

			service := NewTasksService(postgres, redis, log)
			task, err := service.PatchTask(ctx, input.taskID, input.userID, input.command)
			if input.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, input.expErr.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, input.expResult, task)
		})
	}
}
