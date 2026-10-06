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

type testPatchStatusInput struct {
	name      string
	status    string
	mocks     func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task)
	expResult *core_domain.Task
	expErr    error
}

func TestPatchStatus(t *testing.T) {
	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)

	const (
		token          = "40d668c1f2fd4f24a85dfed76d5e76b543540bca17bb928dad4af0bc71671d63"
		invalidStatus  = "creted"
		taskID         = 1
		userID         = 1
		version        = 1
		redisErrorCode = -1
	)

	initialTask := core_domain.NewTask(taskID, version, userID, "Clean the park", "Remove plastic bottles near the park", core_domain.StatusCreated, time.Date(2026, time.September, 28, 18, 0, 0, 0, time.UTC), nil)
	updatedTaskInProgress := core_domain.NewTask(taskID, version, userID, "Clean the park", "Remove plastic bottles near the park", core_domain.StatusInProgress, time.Date(2026, time.September, 28, 18, 0, 0, 0, time.UTC), nil)
	updatedTaskRejected := core_domain.NewTask(taskID, version, userID, "Clean the park", "Remove plastic bottles near the park", core_domain.StatusRejected, time.Date(2026, time.September, 28, 18, 0, 0, 0, time.UTC), nil)

	inputs := []testPatchStatusInput{
		{
			"success: patched status to 'in_progress'",
			string(core_domain.StatusInProgress),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				redis.EXPECT().GetAction(ctx, token).Return(taskID, nil)
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchStatus(ctx, gomock.Any()).Return(updatedTaskInProgress, nil)
				redis.EXPECT().ProlongAction(ctx, token).Return(nil)
			},
			updatedTaskInProgress,
			nil,
		},
		{
			"success: patched status to 'rejected'",
			string(core_domain.StatusRejected),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				redis.EXPECT().GetAction(ctx, token).Return(taskID, nil)
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchStatus(ctx, gomock.Any()).Return(updatedTaskRejected, nil)
				redis.EXPECT().DeleteRedis(ctx, token).Return(nil)
			},
			updatedTaskRejected,
			nil,
		},
		{
			"fail: status validation error",
			invalidStatus,
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
			},
			nil,
			fmt.Errorf("validate status: invalid status: invalid argument"),
		},
		{
			"fail: get action from redis error",
			string(core_domain.StatusInProgress),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				redis.EXPECT().GetAction(ctx, token).Return(redisErrorCode, fmt.Errorf("redis down"))
			},
			nil,
			fmt.Errorf("get task id from redis: redis down"),
		},
		{
			"fail: get task from repository error",
			string(core_domain.StatusInProgress),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				redis.EXPECT().GetAction(ctx, token).Return(taskID, nil)
				repo.EXPECT().GetTask(ctx, taskID).Return(nil, fmt.Errorf("sql: connection refused"))
			},
			nil,
			fmt.Errorf("get task from repository: sql: connection refused"),
		},
		{
			"fail: apply changes error",
			string(core_domain.StatusDone),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				redis.EXPECT().GetAction(ctx, token).Return(taskID, nil)
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
			},
			nil,
			fmt.Errorf("apply status patch: transition from: created to: done is not allowed: invalid argument"),
		},
		{
			"fail: patch status in repository error",
			string(core_domain.StatusInProgress),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				redis.EXPECT().GetAction(ctx, token).Return(taskID, nil)
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchStatus(ctx, updatedTaskInProgress).Return(nil, fmt.Errorf("sql: connection refused"))
			},
			nil,
			fmt.Errorf("change status in repository: sql: connection refused"),
		},
		{
			"fail: prolong action in redis error",
			string(core_domain.StatusInProgress),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				redis.EXPECT().GetAction(ctx, token).Return(taskID, nil)
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchStatus(ctx, updatedTaskInProgress).Return(updatedTaskInProgress, nil)
				redis.EXPECT().ProlongAction(ctx, token).Return(fmt.Errorf("redis down"))
			},
			updatedTaskInProgress,
			nil,
		},
		{
			"fail: delete action in redis error",
			string(core_domain.StatusRejected),
			func(repo *mock_tasks_service.MockTasksRepository, redis *mock_tasks_service.MockTasksRedis, initialTask *core_domain.Task) {
				redis.EXPECT().GetAction(ctx, token).Return(taskID, nil)
				repo.EXPECT().GetTask(ctx, taskID).Return(initialTask, nil)
				repo.EXPECT().PatchStatus(ctx, gomock.Any()).Return(updatedTaskRejected, nil)
				redis.EXPECT().DeleteRedis(ctx, token).Return(fmt.Errorf("redis down"))
			},
			updatedTaskRejected,
			nil,
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
			task, err := service.PatchStatus(ctx, input.status, token)
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
