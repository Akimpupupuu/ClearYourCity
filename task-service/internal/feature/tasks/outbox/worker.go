package tasks_outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	tasks_postgres "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/repository/postgres"
)

const (
	batchTimeout = 5 * time.Second
	dbTimeout    = 1 * time.Second
)

type Worker struct {
	taskRepository TaskRepository
	tasksRedis     TasksRedis
	tasksKafka     TasksKafka
	log            core_logger.Logger
	batchLimit     int
	interval       time.Duration
}

type TaskRepository interface {
	GetPendingMessages(ctx context.Context, limit int) ([]tasks_postgres.MessageModel, error)
	PatchMessageStatus(ctx context.Context, ids []int) error
}

type TasksRedis interface {
	CreateAction(ctx context.Context, token string, taskID int) error
}

type TasksKafka interface {
	PublishTaskCreated(ctx context.Context, taskID int, message []byte) error
}

func NewWorker(
	taskRepository TaskRepository,
	tasksRedis TasksRedis,
	tasksKafka TasksKafka,
	log core_logger.Logger,
	config Config,
) *Worker {
	return &Worker{
		taskRepository: taskRepository,
		tasksRedis:     tasksRedis,
		tasksKafka:     tasksKafka,
		log:            log,
		batchLimit:     config.BatchLimit,
		interval:       config.Interval,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Warn("stopping outbox worker due to context cancellation")
			return
		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				w.log.Error("process outbox batch", core_logger.Err(err))
			}
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) error {
	batchCtx, cancel := context.WithTimeout(ctx, batchTimeout)
	defer cancel()

	messages, err := w.taskRepository.GetPendingMessages(batchCtx, w.batchLimit)
	if err != nil {
		return fmt.Errorf("get messages from repository: %w", err)
	}

	if len(messages) == 0 {
		return nil
	}

	processedIDs := make([]int, 0, len(messages))

	for _, msg := range messages {
		if ctx.Err() != nil {
			break
		}

		var payload struct {
			Token string `json:"token"`
		}

		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			w.log.Error("failed to unmarshal payload", core_logger.Err(err))
			continue
		}

		if err = w.tasksRedis.CreateAction(batchCtx, payload.Token, msg.TaskID); err != nil {
			w.log.Error("failed to create action in redis", core_logger.Err(err))
			break
		}

		if err = w.tasksKafka.PublishTaskCreated(batchCtx, msg.TaskID, msg.Payload); err != nil {
			w.log.Error("failed to produce message into kafka", core_logger.Err(err))
			break
		}

		processedIDs = append(processedIDs, msg.ID)
	}

	if len(processedIDs) > 0 {
		detachedCtx := context.WithoutCancel(ctx)

		updateCtx, updateCancel := context.WithTimeout(detachedCtx, dbTimeout)

		err := w.taskRepository.PatchMessageStatus(updateCtx, processedIDs)
		updateCancel()

		if err != nil {
			return fmt.Errorf("patch status as processed: %w", err)
		}

	}

	return nil
}
