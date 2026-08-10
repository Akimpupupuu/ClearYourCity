package tasks_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	core_kafka "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/kafka"
	tasks_service "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/service"
)

const (
	taskNotificationTopic = "task-created"
)

type tasksKafkaRepository struct {
	producer *core_kafka.Producer
}

func NewTasksKafkaRepository(producer *core_kafka.Producer) *tasksKafkaRepository {
	return &tasksKafkaRepository{
		producer: producer,
	}
}

func (r *tasksKafkaRepository) PublishTaskCreated(ctx context.Context, event tasks_service.ServiceEventDTO) error {
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event data: %w", err)
	}

	key := []byte(strconv.Itoa(event.TaskID))
	if err = r.producer.Produce(ctx, taskNotificationTopic, key, message); err != nil {
		return fmt.Errorf("publish message to kafka: %w", err)
	}

	return nil
}
